package account

import "time"

type Account struct {
	ID            string    `json:"id"`
	CustomerID    string    `json:"customer_id"`
	AccountNumber string    `json:"account_number"`
	Currency      string    `json:"currency"`
	Balance       float64   `json:"balance"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	CustomerID string  `json:"customer_id" binding:"required"`
	Currency   string  `json:"currency" binding:"required,len=3"`
	Balance    float64 `json:"balance" binding:"optional"`
}

type UpdateAccountRequest struct {
	Currency string `json:"currency" binding:"omitempty,len=3"`
	Status   string `json:"status"`
}

type CloseAccountRequest struct {
	Status string `json:"status"`
}
