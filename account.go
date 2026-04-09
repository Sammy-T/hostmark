package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/sammy-t/hostmark/internal/auth"
	"gorm.io/gorm"
)

func handleGetMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		if accessToken == nil {
			http.Error(w, "auth required", http.StatusUnauthorized)
			return
		}

		var user User

		if result := db.Where("username = ?", claims.Subject).First(&user); result.Error != nil {
			msg := "data error"
			code := http.StatusInternalServerError

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				msg = "invalid auth"
				code = http.StatusBadRequest
			}

			http.Error(w, msg, code)
			return
		}

		ruleArgs := auth.RuleArgs{
			User:  user.Username,
			Owner: user.Username,
		}

		if granted := auth.Access(user.Role, auth.ResAcct, auth.PermRead, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResAcct, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		resp, err := json.Marshal(user)
		if err != nil {
			log.Printf("error creating response: %v", err)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}

func handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		if accessToken == nil {
			http.Error(w, "auth required", http.StatusUnauthorized)
			return
		}

		var user User

		if result := db.Where("username = ?", claims.Subject).First(&user); result.Error != nil {
			msg := "data error"
			code := http.StatusInternalServerError

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				msg = "invalid auth"
				code = http.StatusBadRequest
			}

			http.Error(w, msg, code)
			return
		}

		reqUsername := r.PathValue("username")

		ruleArgs := auth.RuleArgs{
			User:  user.Username,
			Owner: reqUsername,
		}

		if granted := auth.Access(user.Role, auth.ResAcct, auth.PermRead, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResAcct, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		var reqUser User

		if result := db.Where("username = ?", reqUsername).First(&reqUser); result.Error != nil {
			msg := "data error"
			code := http.StatusInternalServerError

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				msg = "invalid request"
				code = http.StatusBadRequest
			}

			http.Error(w, msg, code)
			return
		}

		resp, err := json.Marshal(reqUser)
		if err != nil {
			log.Printf("error creating response: %v", err)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}

func handleGetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		if accessToken == nil {
			http.Error(w, "auth required", http.StatusUnauthorized)
			return
		}

		var user User

		if result := db.Where("username = ?", claims.Subject).First(&user); result.Error != nil {
			msg := "data error"
			code := http.StatusInternalServerError

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				msg = "invalid auth"
				code = http.StatusBadRequest
			}

			http.Error(w, msg, code)
			return
		}

		ruleArgs := auth.RuleArgs{
			User: user.Username,
		}

		if granted := auth.Access(user.Role, auth.ResAcct, auth.PermRead, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResAcct, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		var users []User

		if result := db.Find(&users); result.Error != nil {
			log.Print(result.Error)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(users)
		if err != nil {
			log.Printf("json: %v", err)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}

func handleUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("parse form: %v", err)
			http.Error(w, "unable to parse request", http.StatusInternalServerError)
			return
		}

		reqUsername := r.PathValue("username")
		if reqUsername == "" {
			msg := "invalid user"
			log.Print(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		updates := make(map[string]any)
		prefUpdates := make(map[string]any)

		if visibility := r.PostForm.Get("default-visibility"); auth.IsValidVisibility(visibility) {
			prefUpdates["note_vis"] = visibility
		}

		if len(updates) == 0 && len(prefUpdates) == 0 {
			msg := "invalid fields"
			log.Print(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		if accessToken == nil {
			http.Error(w, "auth required", http.StatusUnauthorized)
			return
		}

		var user User

		if result := db.Where("username = ?", claims.Subject).First(&user); result.Error != nil {
			msg := "data error"
			code := http.StatusInternalServerError

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				msg = "invalid auth"
				code = http.StatusBadRequest
			}

			http.Error(w, msg, code)
			return
		}

		ruleArgs := auth.RuleArgs{
			User:  user.Username,
			Owner: reqUsername,
		}

		if granted := auth.Access(user.Role, auth.ResAcct, auth.PermUpdate, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResAcct, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		if len(prefUpdates) > 0 {
			if result := db.Model(&Preferences{}).Where("user = ?", reqUsername).Updates(prefUpdates); result.Error != nil {
				http.Error(w, "unable to update user", http.StatusInternalServerError)
				return
			}
		}

		if len(updates) > 0 {
			if result := db.Model(&User{}).Where("username = ?", reqUsername).Updates(updates); result.Error != nil {
				http.Error(w, "unable to update user", http.StatusInternalServerError)
				return
			}
		}
	}
}
