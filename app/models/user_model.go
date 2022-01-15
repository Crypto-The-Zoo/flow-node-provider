package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

type UserPublic struct {
	ID          uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Username    string    `db:"username" json:"username" validate:"required,lte=255"`
	Email       string    `db:"email" json:"email" validate:"required,lte=255"`
	FlowAddress string    `db:"flow_address" json:"flow_address" validate:"lte=255"`
}

type LoginObj struct {
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Value make the LoginObj struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (l LoginObj) Value() (driver.Value, error) {
	return json.Marshal(l)
}

// Scan make the LoginObj struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (l *LoginObj) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &l)
}
