package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sammy-t/hostmark/internal/auth"
	"github.com/sammy-t/hostmark/pwd"
	"golang.org/x/crypto/argon2"
)

var hmSecret string = "my-hostmark-secret" //// TODO: Load from env

var hashParams pwd.HashParams = pwd.HashParams{
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

func handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %q", r.Method, r.URL.String())

		err := r.ParseForm()
		if err != nil {
			log.Printf("parse form: %v", err)
			http.Error(w, "unable to parse request", http.StatusInternalServerError)
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		if !auth.IsValidUsername(username) {
			err = fmt.Errorf("invalid username")
		}

		//// TODO: TEMP
		// if !auth.IsValidPassword(password) {
		// 	err = fmt.Errorf("invalid password")
		// }

		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//// TODO: TEMP
		// if err = pwd.CheckAgainstPwned("hostmark.sammy-t", password, 25); err != nil {
		// 	log.Print(err)
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		foundResult := db.Where("username = ?", username).Limit(1).Find(&User{})

		if foundResult.Error != nil {
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
		}
		log.Print(user)

		//// TODO: TEMP
		// if result := db.Create(&user); result.Error != nil {
		// 	http.Error(w, "unable to create user", http.StatusInternalServerError)
		// 	return
		// }

		//// TODO: Set reasonable durations
		accessDuration := 2 * time.Minute
		refreshDuration := 5 * time.Minute

		refreshId, err := uuid.NewV7()
		if err != nil {
			log.Printf("uuid: %v", err)
			http.Error(w, "error completing auth", http.StatusInternalServerError)
			return
		}

		nonceBytes := pwd.GenerateRandBytes(32)
		nonce := base64.StdEncoding.EncodeToString(nonceBytes)

		accessClaims := jwt.MapClaims{
			"iss": "hostmark",
			"aud": "acc",
			"sub": username,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(accessDuration).Unix(),
		}

		refreshClaims := jwt.MapClaims{
			"iss": "hostmark",
			"aud": "ref",
			"sub": username,
			"jti": refreshId.String(),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(refreshDuration).Unix(),
		}

		deviceClaims := jwt.MapClaims{
			"iss": "hostmark",
			"aud": "dev",
			"sub": username,
			"jti": nonce,
		}

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
		accessJwt, err := accessToken.SignedString([]byte(hmSecret))
		if err != nil {
			msg := "error completing auth"
			log.Printf("%v: %v", msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
		refreshJwt, err := refreshToken.SignedString([]byte(hmSecret))
		if err != nil {
			msg := "error completing auth"
			log.Printf("%v: %v", msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		deviceToken := jwt.NewWithClaims(jwt.SigningMethodHS256, deviceClaims)
		deviceJwt, err := deviceToken.SignedString([]byte(hmSecret))
		if err != nil {
			msg := "error completing auth"
			log.Printf("%v: %v", msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		accessCookie := http.Cookie{
			Name:     "client_acc",
			Value:    accessJwt,
			Path:     "/api",
			MaxAge:   int(accessDuration.Seconds()),
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
			Secure:   true,
		}

		refreshCookie := http.Cookie{
			Name:     "client_ref",
			Value:    refreshJwt,
			Path:     "/api",
			MaxAge:   int(refreshDuration.Seconds()),
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
			Secure:   true,
		}

		deviceCookie := http.Cookie{
			Name:     "client_dev",
			Value:    deviceJwt,
			Path:     "/api",
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
			Secure:   true,
		}

		http.SetCookie(w, &accessCookie)
		http.SetCookie(w, &refreshCookie)
		http.SetCookie(w, &deviceCookie)
	}
}
