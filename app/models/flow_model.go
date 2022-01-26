package models

import "time"

type Block struct {
	ID        string
	Height    string
	Timestamp time.Time
}
