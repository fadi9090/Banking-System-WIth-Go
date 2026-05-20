package address

import "time"

type Address struct {
	ID         string    `json:"id"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	ZipCode    string    `json:"zip_code"`
	Country    string    `json:"country"`
	CustomerID string    `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateAddressRequest struct {
	Street     string `json:"street" binding:"required"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state" binding:"required"`
	ZipCode    string `json:"zip_code" binding:"required"`
	Country    string `json:"country" binding:"required"`
	CustomerID string `json:"customer_id" binding:"required"`
}

type UpdateAddressRequest struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}
