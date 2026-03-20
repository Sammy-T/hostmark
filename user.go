package main

import dbExt "github.com/sammy-t/hostmark/db"

type User struct {
	dbExt.Model
	Username string `gorm:"unique"`
	PwdHash  string
}
