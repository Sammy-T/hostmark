package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/sammy-t/hostmark/internal/auth"
	"github.com/sammy-t/hostmark/pwd"
)

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
			http.Error(w, "unable to parse request", 500)
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
			http.Error(w, err.Error(), 400)
			return
		}

		//// TODO: TEMP
		// if err = pwd.CheckAgainstPwned("hostmark.sammy-t", password, 25); err != nil {
		// 	log.Print(err)
		// 	http.Error(w, err.Error(), 400)
		// 	return
		// }

		foundResult := db.Where("username = ?", username).Limit(1).Find(&User{})

		if foundResult.Error != nil {
			http.Error(w, "data error", 500)
			return
		} else if foundResult.RowsAffected != 0 {
			log.Printf("username %q already exists", username)
			http.Error(w, "invalid username", 400)
			return
		}

		s := pwd.GenerateSalt(32)
		h := pwd.HashPwd([]byte(password), s, hashParams)

		salt := base64.StdEncoding.EncodeToString(s)
		hashed := base64.StdEncoding.EncodeToString(h)

		user := User{
			Username: username,
			PwdHash:  hashed,
			Salt:     salt,
		}

		if result := db.Create(&user); result.Error != nil {
			http.Error(w, "unable to create user", 500)
			return
		}

		//// TODO: Return access token, refresh token, and device cookies
	}
}
