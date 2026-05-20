package card

import "time"

type Card struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Last4     string    `json:"last4"`
	Brand     string    `json:"brand"`
	IsBlocked bool      `json:"is_blocked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCardRequest struct {
	AccountID string `json:"account_id" binding:"required"`
	Brand     string `json:"brand" binding:"required"`
}

type UpdateCardRequest struct {
	Brand     string `json:"brand"`
	IsBlocked *bool  `json:"is_blocked"`
}

type BlockCardRequest struct {
	IsBlocked bool `json:"is_blocked"`
}
