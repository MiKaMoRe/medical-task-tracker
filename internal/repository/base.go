package repository

import "time"

type Base struct {
	ID        int       `gorm:"primarykey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
}
