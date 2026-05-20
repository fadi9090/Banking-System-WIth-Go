package customer

import "time"

type Customer struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PersonID       string    `json:"person_id"`
	CustomerNumber string    `json:"customer_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateCustomerRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	PersonID string `json:"person_id" binding:"required"`
}

type UpdateCustomerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"omitempty,email"`
}
