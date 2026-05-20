package transaction

import "time"

type Transaction struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	CardID      *string   `json:"card_id,omitempty"`
	Direction   string    `json:"direction"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Reference   string    `json:"reference"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateTransactionRequest struct {
	AccountID   string  `json:"account_id" binding:"required"`
	CardID      *string `json:"card_id"`
	Direction   string  `json:"direction" binding:"required,oneof=debit credit"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3"`
	Description string  `json:"description"`
}

type UpdateTransactionRequest struct {
	Status string `json:"status"`
}
