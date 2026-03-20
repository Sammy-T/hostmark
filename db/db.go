package db

import "time"

// Like gorm.Model except without the `DeletedAt` field
// Gorm uses for soft deletion.
//
// See: https://gorm.io/docs/models.html#gorm-Model
type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
