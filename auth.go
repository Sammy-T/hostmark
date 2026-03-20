package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sammy-t/hostmark/internal/auth"
	"github.com/sammy-t/hostmark/pwd"
)

func handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %q", r.Method, r.URL.String())

		err := r.ParseForm()
		if err != nil {
			log.Printf("parse form: %v", err)
			http.Error(w, "unable to parse request", 500)
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		if !auth.IsValidUsername(username) {
			err = fmt.Errorf("invalid username")
		}

		if !auth.IsValidPassword(password) {
			err = fmt.Errorf("invalid password")
		}

		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), 400)
			return
		}

		if err = pwd.CheckAgainstPwned("hostmark.sammy-t", password, 25); err != nil {
			log.Print(err)
			http.Error(w, err.Error(), 400)
			return
		}

		//// TODO: Check if username already exists
		//// TODO: Create user
		//// TODO: Return access token, refresh token, and device cookies
	}
}
