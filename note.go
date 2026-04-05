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

func handleGetTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, _ := parseToken(CookieAccess, accessCookie)

		if accessToken == nil {
			http.Error(w, "auth required", http.StatusUnauthorized)
			return
		}

		var tags []Tag

		if result := db.Find(&tags); result.Error != nil {
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(tags)
		if err != nil {
			log.Printf("error creating response: %v", err)
			http.Error(w, "data error", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}

func handleGetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessCookie, _ := r.Cookie(string(CookieAccess))
		accessToken, claims := parseToken(CookieAccess, accessCookie)

		var user User

		if accessToken != nil {
			db.Where("username = ?", claims.Subject).First(&user)
		}

		var notes []Note
		var authCond *gorm.DB // The query condition corresponding to the request's auth.

		// Adjust the database query according to whether the request has auth
		if accessToken != nil && user.Username != "" {
			authCond = db.Where("owner = ? OR visibility IN ?", user.Username, []string{"protected", "public"})
		} else {
			authCond = db.Where("visibility = ?", "public")
		}

		tags := r.URL.Query()["tags"]

		if len(tags) > 0 {
			// Gorm doesn't have utilities for querying/filtering by associations.
			//
			// See: https://github.com/go-gorm/gorm/issues/3287#issuecomment-908893840
			db.Preload(clause.Associations).
				Joins("JOIN note_tags ON note_tags.note_id = notes.id").
				Joins("JOIN tags ON note_tags.tag_name = tags.name").
				Where("tags.name in ?", tags).
				Where(authCond).
				Order("notes.created_at DESC").
				Find(&notes)
		} else {
			db.Preload(clause.Associations).
				Where(authCond).
				Order("created_at DESC").
				Find(&notes)
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

		visibility := r.PostForm.Get("visibility")
		content := r.PostForm.Get("content")

		if content == "" || !auth.IsValidVisibility(visibility) {
			log.Print("no fields to create")
			http.Error(w, "invalid fields", http.StatusBadRequest)
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

		tagNames := r.PostForm["tags"]

		var tags []*Tag

		for _, name := range tagNames {
			tags = append(tags, &Tag{Name: name})
		}

		note := Note{
			Owner:      user.Username,
			Visibility: visibility,
			Content:    content,
		}

		if len(tags) > 0 {
			note.Tags = tags
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

func handleUpdateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			msg := "invalid note id"
			log.Print(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		if err = r.ParseForm(); err != nil {
			log.Printf("parse form: %v", err)
			http.Error(w, "unable to parse request", http.StatusInternalServerError)
			return
		}

		updates := make(map[string]any)

		if visibility := r.PostForm.Get("visibility"); auth.IsValidVisibility(visibility) {
			updates["visibility"] = visibility
		}

		if content := r.PostForm.Get("content"); content != "" {
			updates["content"] = content
		}

		var tagNames []string

		if tagStr := r.PostForm.Get("tags"); tagStr != "" {
			tagNames = strings.Split(tagStr, ",")
		}

		if len(updates) == 0 && len(tagNames) == 0 {
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

			log.Printf("find user: %v", result.Error)
			http.Error(w, msg, code)
			return
		}

		var note Note

		if result := db.Where("id = ?", id).First(&note); result.Error != nil {
			http.Error(w, "invalid note id", http.StatusBadRequest)
			return
		}

		ruleArgs := auth.RuleArgs{
			User:  user.Username,
			Owner: note.Owner,
		}

		if granted := auth.Access(user.Role, auth.ResNote, auth.PermUpdate, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResNote, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		if len(updates) > 0 {
			if result := db.Model(&note).Updates(updates); result.Error != nil {
				http.Error(w, "unable to update note", http.StatusInternalServerError)
				return
			}
		}

		if len(tagNames) > 0 {
			var tags []*Tag

			for _, name := range tagNames {
				tags = append(tags, &Tag{Name: name})
			}

			if err = db.Model(&note).Association("Tags").Replace(tags); err != nil {
				http.Error(w, "unable to update tags", http.StatusInternalServerError)
				return
			}
		}
	}
}

func handleDelNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			msg := "invalid note id"
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

			log.Printf("find user: %v", result.Error)
			http.Error(w, msg, code)
			return
		}

		var note Note

		if result := db.Where("id = ?", id).First(&note); result.Error != nil {
			http.Error(w, "invalid note id", http.StatusBadRequest)
			return
		}

		ruleArgs := auth.RuleArgs{
			User:  user.Username,
			Owner: note.Owner,
		}

		if granted := auth.Access(user.Role, auth.ResNote, auth.PermDelete, ruleArgs); !granted {
			log.Printf("access denied for %v to %v", auth.ResNote, user.Username)
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		if result := db.Delete(&note); result.Error != nil {
			http.Error(w, "unable to delete note", http.StatusInternalServerError)
			return
		}
	}
}
