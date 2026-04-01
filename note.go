package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sammy-t/hostmark/internal/auth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func handleGetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		var user User

		if accessToken != nil {
			db.Where("username = ?", claims.Subject).First(&user)
		}

		var notes []Note

		// Adjust the database query according to whether the request has auth
		if accessToken != nil && user.Username != "" {
			db.Where("owner = ? OR visibility IN ?", user.Username, []string{"protected", "public"}).Preload(clause.Associations).Find(&notes)
		} else {
			db.Where("visibility = ?", "public").Preload(clause.Associations).Find(&notes)
		}

		resp, err := json.Marshal(notes)
		if err != nil {
			log.Printf("error creating response: %v", err)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}

func handleCreateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("parse form: %v", err)
			http.Error(w, "unable to parse request", http.StatusInternalServerError)
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

			log.Printf("find user: %v", result.Error)
			http.Error(w, msg, code)
			return
		}

		tagNames := strings.Split(r.PostForm.Get("tags"), ",")

		var tags []*Tag

		for _, name := range tagNames {
			tags = append(tags, &Tag{Name: name})
		}

		note := Note{
			Owner:      user.Username,
			Visibility: r.PostForm.Get("visibility"),
			Tags:       tags,
			Content:    r.PostForm.Get("content"),
		}

		if note.Content == "" {
			msg := "no note content"
			log.Print(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		ruleArgs := auth.RuleArgs{
			User:       user.Username,
			Owner:      note.Owner,
			Visibility: note.Visibility,
		}

		if granted := auth.Access(user.Role, auth.ResNote, auth.PermCreate, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResNote, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		if result := db.Create(&note); result.Error != nil {
			http.Error(w, "unable to create note", http.StatusInternalServerError)
			return
		}
	}
}

func handleGetNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		var user User

		if accessToken != nil {
			db.Where("username = ?", claims.Subject).First(&user)
		}

		// The path value string could likely be used without converting
		// but doing so validates that the provided value is an int.
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			msg := "invalid note id"
			log.Print(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		var note Note

		if result := db.Where("id = ?", id).Preload(clause.Associations).First(&note); result.Error != nil {
			http.Error(w, "invalid note id", http.StatusBadRequest)
			return
		}

		ruleArgs := auth.RuleArgs{
			User:       user.Username,
			Owner:      note.Owner,
			Visibility: note.Visibility,
		}

		if granted := auth.Access(user.Role, auth.ResNote, auth.PermRead, ruleArgs); !granted {
			name := "unknown"
			if user.Username != "" {
				name = user.Username
			}

			log.Printf("access denied for %v to %v", auth.ResNote, name)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		resp, err := json.Marshal(note)
		if err != nil {
			log.Printf("error creating response: %v", err)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}
