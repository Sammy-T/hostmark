package main

import (
	"time"

	dbExt "github.com/sammy-t/hostmark/db"
	"github.com/sammy-t/hostmark/internal/auth"
)

type User struct {
	dbExt.Model
	Username     string     `gorm:"unique" json:"username"`
	PwdHash      string     `json:"-"`
	Salt         string     `json:"-"`
	LockdownTime *time.Time `json:"-"`
	Role         auth.Role  `gorm:"default:user" json:"role"`
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
