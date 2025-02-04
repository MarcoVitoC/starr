package internal

import (
	"time"

	"github.com/google/uuid"
)

type Wish struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
