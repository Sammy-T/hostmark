package main

import (
	"time"

	dbExt "github.com/sammy-t/hostmark/db"
)

type User struct {
	dbExt.Model
	Username     string `gorm:"unique"`
	PwdHash      string
	Salt         string
	LockdownTime *time.Time
}

type FailedLogin struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	Username  string
	Nonce     *string
}

type LockedToken struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	Username  string
	Nonce     string `gorm:"unique"`
}

type RefreshToken struct {
	ID        string `gorm:"primaryKey"`
	Username  string
	IssuedAt  time.Time
	ExpiresAt time.Time
}
