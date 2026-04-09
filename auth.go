package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sammy-t/hostmark/internal/auth"
	"github.com/sammy-t/hostmark/pwd"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type CookieName string

const (
	CookieAccess  CookieName = "client_acc"
	CookieRefresh CookieName = "client_ref"
	CookieDevice  CookieName = "client_dev"
)

var hmSecret string = "my-hostmark-secret" //// TODO: Load from env

var accessDuration time.Duration = 2 * time.Minute
var refreshDuration time.Duration = 5 * time.Minute
var lockDuration time.Duration = 5 * time.Minute
var authAttemptLimit int = 3

var hashParams pwd.HashParams = pwd.HashParams{
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

func handleSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("parse form: %v", err)
			http.Error(w, "unable to parse request", http.StatusInternalServerError)
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		if !auth.IsValidPassword(password) {
			err = fmt.Errorf("invalid password")
		}

		if !auth.IsValidUsername(username) {
			err = fmt.Errorf("invalid username")
		}

		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = pwd.CheckAgainstPwned("hostmark.sammy-t", password, 25); err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		foundResult := db.Where("username = ?", username).Limit(1).Find(&User{})

		if foundResult.Error != nil {
			log.Printf("find user: %v", foundResult.Error)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		} else if foundResult.RowsAffected != 0 {
			log.Printf("username %q already exists", username)
			http.Error(w, "invalid username", http.StatusBadRequest)
			return
		}

		s := pwd.GenerateRandBytes(32)
		h := argon2.IDKey([]byte(password), s, hashParams.Time, hashParams.Memory, hashParams.Threads, hashParams.KeyLen)

		salt := base64.StdEncoding.EncodeToString(s)
		hashed := base64.StdEncoding.EncodeToString(h)

		user := User{
			Username: username,
			PwdHash:  hashed,
			Salt:     salt,
			Prefs:    Preferences{NoteVis: "private"},
		}

		if result := db.Create(&user); result.Error != nil {
			http.Error(w, "unable to create user", http.StatusInternalServerError)
			return
		}

		errMsg := "error completing auth"

		accessCookie, _, _, err := createTokenCookie(CookieAccess, username)
		if err != nil {
			log.Print(err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		refreshCookie, refreshToken, _, err := createTokenCookie(CookieRefresh, username)
		if err != nil {
			log.Print(err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		deviceCookie, _, _, err := createTokenCookie(CookieDevice, username)
		if err != nil {
			log.Print(err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		if err = storeRefreshToken(refreshToken); err != nil {
			log.Printf("store ref: %v", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(user)
		if err != nil {
			log.Printf("error creating response: %v", err)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, accessCookie)
		http.SetCookie(w, refreshCookie)
		http.SetCookie(w, deviceCookie)

		w.Write(resp)
	}
}

func handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("parse form: %v", err)
			http.Error(w, "unable to parse request", http.StatusInternalServerError)
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		if !auth.IsValidPassword(password) {
			err = fmt.Errorf("invalid password")
		}

		if !auth.IsValidUsername(username) {
			err = fmt.Errorf("invalid username")
		}

		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		reqDeviceCookie, _ := r.Cookie(string(CookieDevice))
		reqDeviceToken, claims := parseToken(CookieDevice, reqDeviceCookie)

		var nonce string

		if reqDeviceToken != nil && claims.Subject == username && claims.ID != "" { // Trusted client
			nonce = claims.ID

			var lockedTokens []LockedToken

			if result := db.Where("username = ?", username).Find(&lockedTokens); result.Error != nil {
				log.Printf("locked tokens: %v", result.Error)
				http.Error(w, "data error", http.StatusInternalServerError)
				return
			}

			var expiredLocks []LockedToken
			var locked bool

			for _, token := range lockedTokens {
				expired := time.Now().After(token.CreatedAt.Add(lockDuration))

				if expired {
					expiredLocks = append(expiredLocks, token)
				} else if token.Nonce == nonce {
					locked = true
				}
			}

			if len(expiredLocks) > 0 {
				db.Delete(&expiredLocks)
			}

			if locked {
				log.Print("auth attempt with locked device token")
				http.Error(w, "invalid credentials", http.StatusBadRequest)
				return
			}
		} else { // Untrusted client
			var user User

			if result := db.Where("username = ?", username).First(&user); result.Error != nil {
				msg := "data error"
				code := http.StatusInternalServerError

				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					msg = "invalid credentials"
					code = http.StatusBadRequest
				}

				log.Printf("find user: %v", result.Error)
				http.Error(w, msg, code)
				return
			}

			if user.LockdownTime != nil {
				expired := time.Now().After(user.LockdownTime.Add(lockDuration))

				if expired {
					log.Print("removing account lockdown")

					if result := db.Model(&user).Update("lockdown_time", nil); result.Error != nil {
						log.Printf("remove lockdown: %v", result.Error)
						http.Error(w, "data error", http.StatusInternalServerError)
						return
					}
				} else {
					log.Print("auth attempt with locked account")
					http.Error(w, "invalid credentials", http.StatusBadRequest)
					return
				}
			}
		}

		var user User

		if result := db.Preload("Prefs").Where("username = ?", username).First(&user); result.Error != nil {
			msg := "data error"
			code := http.StatusInternalServerError

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				msg = "invalid credentials"
				code = http.StatusBadRequest
			}

			log.Printf("find user: %v", result.Error)
			http.Error(w, msg, code)
			return
		}

		s, err := base64.StdEncoding.DecodeString(user.Salt)
		if err != nil {
			log.Print(err)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		h := argon2.IDKey([]byte(password), s, hashParams.Time, hashParams.Memory, hashParams.Threads, hashParams.KeyLen)

		hashed := base64.StdEncoding.EncodeToString(h)

		if hashed != user.PwdHash {
			log.Print("invalid password")

			failed := FailedLogin{
				Username: username,
			}

			if reqDeviceToken != nil && nonce != "" {
				failed.Nonce = &nonce
			}

			if result := db.Create(&failed); result.Error != nil {
				log.Printf("failed login: %v", result.Error)
				http.Error(w, "data error", http.StatusInternalServerError)
				return
			}

			var stored []FailedLogin
			var loginsResult *gorm.DB

			if failed.Nonce != nil {
				loginsResult = db.Where("username = ? AND nonce = ?", username, failed.Nonce).Find(&stored)
			} else {
				loginsResult = db.Where("username = ?", username).Find(&stored)
			}

			if loginsResult.Error != nil {
				log.Printf("failed logins: %v", loginsResult.Error)
				http.Error(w, "data error", http.StatusInternalServerError)
				return
			}

			var expired []FailedLogin
			var failedWithinDur []FailedLogin

			for _, attempt := range stored {
				exp := time.Now().After(attempt.CreatedAt.Add(lockDuration))

				if exp {
					expired = append(expired, attempt)
				} else {
					failedWithinDur = append(failedWithinDur, attempt)
				}
			}

			if len(expired) > 0 {
				db.Delete(&expired)
			}

			log.Printf("found %d, %d/%d failed", len(stored), len(failedWithinDur), authAttemptLimit)

			if len(failedWithinDur) > authAttemptLimit {
				if failed.Nonce != nil {
					locked := LockedToken{
						Username: username,
						Nonce:    *failed.Nonce,
					}

					log.Print("locking token")

					if result := db.Create(&locked); result.Error != nil {
						log.Printf("locked token: %v", result.Error)
						http.Error(w, "data error", http.StatusInternalServerError)
						return
					}
				} else {
					log.Print("lockdown account")

					if result := db.Model(&user).Update("lockdown_time", time.Now()); result.Error != nil {
						log.Printf("set lockdown: %v", result.Error)
						http.Error(w, "data error", http.StatusInternalServerError)
						return
					}
				}
			}

			http.Error(w, "invalid credentials", http.StatusBadRequest)
			return
		}

		errMsg := "error completing auth"

		accessCookie, _, _, err := createTokenCookie(CookieAccess, username)
		if err != nil {
			log.Print(err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		refreshCookie, refreshToken, _, err := createTokenCookie(CookieRefresh, username)
		if err != nil {
			log.Print(err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		deviceCookie, _, _, err := createTokenCookie(CookieDevice, username)
		if err != nil {
			log.Print(err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		if err = storeRefreshToken(refreshToken); err != nil {
			log.Printf("store ref: %v", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(user)
		if err != nil {
			log.Printf("error creating response: %v", err)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, accessCookie)
		http.SetCookie(w, refreshCookie)
		http.SetCookie(w, deviceCookie)

		w.Write(resp)
	}
}

func handleRefresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqRefreshCookie, _ := r.Cookie(string(CookieRefresh))
		reqRefreshToken, claims := parseToken(CookieRefresh, reqRefreshCookie)

		if reqRefreshToken == nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		username := claims.Subject
		refreshId := claims.ID
		issuedAt := claims.IssuedAt.Time
		expiresAt := claims.ExpiresAt.Time

		if result := db.Where("expires_at < ?", time.Now()).Delete(&RefreshToken{}); result.RowsAffected > 0 {
			log.Printf("deleted %d expired refresh tokens", result.RowsAffected)
		}

		var stored RefreshToken

		query := "id = ? AND username = ?"

		if result := db.First(&stored, query, refreshId, username); result.Error != nil {
			log.Print("token not found")

			if res := db.Where("username = ?", username).Delete(&RefreshToken{}); res.Error == nil {
				log.Printf("invalidated all %d refresh tokens for %v", res.RowsAffected, username)
			}

			log.Print(result.Error)
			http.Error(w, "invalid auth", http.StatusUnauthorized)
			return
		}

		if !issuedAt.Equal(stored.IssuedAt) || !expiresAt.Equal(stored.ExpiresAt) {
			log.Print("token mismatch")

			if res := db.Where("username = ?", username).Delete(&RefreshToken{}); res.Error == nil {
				log.Printf("invalidated all %d refresh tokens for %v", res.RowsAffected, username)
			}

			http.Error(w, "invalid auth", http.StatusUnauthorized)
			return
		}

		if result := db.Delete(&stored); result.Error == nil {
			log.Printf("used %v", stored)
		}

		errMsg := "error completing auth"

		accessCookie, _, _, err := createTokenCookie(CookieAccess, username)
		if err != nil {
			log.Print(err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		refreshCookie, refreshToken, _, err := createTokenCookie(CookieRefresh, username)
		if err != nil {
			log.Print(err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		if err = storeRefreshToken(refreshToken); err != nil {
			log.Printf("store ref: %v", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, accessCookie)
		http.SetCookie(w, refreshCookie)
	}
}

func handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unsetAccess := &http.Cookie{
			Name:   string(CookieAccess),
			Path:   "/api",
			MaxAge: 0,
		}

		unsetRefresh := &http.Cookie{
			Name:   string(CookieRefresh),
			Path:   "/api/auth/refresh",
			MaxAge: 0,
		}

		http.SetCookie(w, unsetAccess)
		http.SetCookie(w, unsetRefresh)
	}
}

func createTokenCookie(name CookieName, username string) (cookie *http.Cookie, token *jwt.Token, tokenJwt string, err error) {
	var claims jwt.RegisteredClaims

	switch name {
	case CookieAccess:
		claims = jwt.RegisteredClaims{
			Issuer:    "hostmark",
			Audience:  jwt.ClaimStrings{"acc"},
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessDuration).UTC()),
		}
	case CookieRefresh:
		refreshId, err := uuid.NewV7()
		if err != nil {
			return nil, nil, "", err
		}

		claims = jwt.RegisteredClaims{
			Issuer:    "hostmark",
			Audience:  jwt.ClaimStrings{"ref"},
			Subject:   username,
			ID:        refreshId.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshDuration).UTC()),
		}
	case CookieDevice:
		nonceBytes := pwd.GenerateRandBytes(32)
		nonce := base64.StdEncoding.EncodeToString(nonceBytes)

		claims = jwt.RegisteredClaims{
			Issuer:   "hostmark",
			Audience: jwt.ClaimStrings{"dev"},
			Subject:  username,
			ID:       nonce,
		}
	default:
		return nil, nil, "", fmt.Errorf("invalid cookie name")
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenJwt, err = token.SignedString([]byte(hmSecret))

	if err != nil {
		return nil, nil, "", err
	}

	cookie = &http.Cookie{
		Name:     string(name),
		Value:    tokenJwt,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}

	switch name {
	case CookieAccess:
		cookie.Path = "/api"
		cookie.MaxAge = int(accessDuration.Seconds())
	case CookieRefresh:
		cookie.Path = "/api/auth/refresh"
		cookie.MaxAge = int(refreshDuration.Seconds())
	case CookieDevice:
		cookie.Path = "/api/auth/login"
	}

	return
}

func parseToken(name CookieName, cookie *http.Cookie) (*jwt.Token, *jwt.RegisteredClaims) {
	if cookie == nil || cookie.Value == "" {
		log.Print("no token cookie")
		return nil, nil
	}

	keyfunc := func(t *jwt.Token) (any, error) {
		return []byte(hmSecret), nil
	}

	opts := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("hostmark"),
	}

	switch name {
	case CookieAccess:
		addl := []jwt.ParserOption{
			jwt.WithAudience("acc"),
			jwt.WithIssuedAt(),
			jwt.WithExpirationRequired(),
		}

		opts = append(opts, addl...)
	case CookieRefresh:
		addl := []jwt.ParserOption{
			jwt.WithAudience("ref"),
			jwt.WithIssuedAt(),
			jwt.WithExpirationRequired(),
		}

		opts = append(opts, addl...)
	case CookieDevice:
		addl := []jwt.ParserOption{
			jwt.WithAudience("dev"),
		}

		opts = append(opts, addl...)
	default:
		log.Print("invalid cookie name")
		return nil, nil
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.RegisteredClaims{}, keyfunc, opts...)
	if err != nil {
		log.Print(err)
		return nil, nil
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return token, claims
	} else {
		log.Print("invalid claims type")
		return nil, nil
	}
}

func storeRefreshToken(token *jwt.Token) error {
	claims, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok {
		return fmt.Errorf("error casting to RegisteredClaims")
	}

	refreshId := claims.ID
	username := claims.Subject

	// Convert the token's UTC time to the local time used by the server.
	refreshIat := claims.IssuedAt.Time.Local()
	refreshExp := claims.ExpiresAt.Time.Local()

	storeToken := RefreshToken{
		ID:        refreshId,
		Username:  username,
		IssuedAt:  refreshIat,
		ExpiresAt: refreshExp,
	}

	if result := db.Create(&storeToken); result.Error != nil {
		return result.Error
	}

	return nil
}
