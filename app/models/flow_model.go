package models

import (
	"time"

	"github.com/onflow/flow-go-sdk"
)

type Block struct {
	ID        string
	Height    string
	Timestamp time.Time
}

type BlockEventsWithHash struct {
	BlockID        flow.Identifier
	Height         uint64
	BlockTimestamp time.Time
	Events         []flow.Event
	BlockIDHash    string
}

// BlockEvents are the events that occurred in a specific block.
type BlockEvents struct {
	ID                string            `json:"id"`
	BlockEventData    map[string]string `json:"blockEventData"`
	EventDate         time.Time         `json:"eventDate"`
	FlowTransactionID string            `json:"flowTransactionId"`

	BlockID          string            `json:"blockId"`
	BlockHeight      uint64            `json:"blockHeight"`
	BlockTimestamp   time.Time         `json:"blockTimestamp"`
	Type             string            `json:"type"`
	TransactionID    string            `json:"transactionId"`
	TransactionIndex int               `json:"transactionIndex"`
	EventIndex       int               `json:"eventIndex"`
	Data             map[string]string `json:"data"`
}

type BlockEventData struct {
	ID         string `json:"id"`
	TemplateID string `json:"templateID"`
}

type Event struct{}

type NFTTemplate struct {
	TypeID      int        `json:"typeId" validate:"required"`
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description" validate:"required"`
	MintLimit   int        `json:"mintLimit" validate:"required"`
	PriceUSD    string     `json:"priceUsd" validate:"required"`
	PriceFlow   string     `json:"priceFlow" validate:"required"`
	Metadata    Metadata   `json:"metadata" validate:"required"`
	Timestamps  Timestamps `json:"timestamps" validate:"required"`
	IsPack      bool       `json:"isPack"`
	IsLand      bool       `json:"isLand"`
}

type Metadata struct {
	Uri      string `json:"uri" validate:"required"`
	Mimetype string `json:"mimetype" validate:"required"`
	Quality  string `json:"quality" validate:"required"`
}

type Timestamps struct {
	AvailableAt time.Time `json:"availableAt" validate:"required"`
	ExpiresAt   time.Time `json:"expiresAt" validate:"required"`
}

type MintNFTRequest struct {
	TypeID  int    `json:"typeID" validate:"required"`
	Address string `json:"Address" validate:"required"`
}
