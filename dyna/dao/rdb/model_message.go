package rdb

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`

	Contents string `json:"content"`
	Name     string `json:"name"`
}
