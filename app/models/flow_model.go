package models

import "time"

type Block struct {
	ID        string
	Height    string
	Timestamp time.Time
}

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
