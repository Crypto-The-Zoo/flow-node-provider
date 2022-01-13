package models

import (
	"time"

	"github.com/google/uuid"
)

// User struct to describe user object.
type User struct {
	ID          uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Username    string    `db:"username" json:"username" validate:"required,lte=255"`
	Email       string    `db:"email" json:"email" validate:"required,lte=255"`
	FlowAddress string    `db:"flow_address" json:"flow_address" validate:"lte=255"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	LoginObj    LoginObj  `db:"login_obj" json:"login_obj"`
}

type LoginObj struct {
	Code      string `json:"code"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}
