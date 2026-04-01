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
	Notes        []Note     `gorm:"foreignKey:Owner;references:Username;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"notes"`
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

type Tag struct {
	Name      string    `gorm:"primaryKey" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Notes     []*Note   `gorm:"many2many:note_tags;" json:"notes"`
}

type Note struct {
	dbExt.Model
	Owner      string `json:"owner"`
	Visibility string `gorm:"default:private" json:"visibility"`
	Tags       []*Tag `gorm:"many2many:note_tags;" json:"tags"`
	Content    string `json:"content"`
}
