package account

import (
	"database/sql"
	"errors"
	"fmt"
	"task-project/internal/components/card"
	"task-project/internal/utility"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(a *CreateAccountRequest) (*Account, error) {
	// Check if customer exists
	cus_query := `SELECT id FROM customers WHERE id = $1`
	var customerID string

	err := r.db.QueryRow(cus_query, a.CustomerID).Scan(&customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer with ID %v does not exist", a.CustomerID)
		}
		return nil, fmt.Errorf("error checking customer: %w", err)
	}

	// Generate account number
	accountNumber := utility.GenerateAccountNumber()

	// Insert the account
	query := `
        INSERT INTO accounts (customer_id, account_number, currency, balance, status, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `

	account := &Account{
		CustomerID:    a.CustomerID,
		AccountNumber: accountNumber,
		Currency:      a.Currency,
		Balance:       a.Balance,
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Execute and scan the returned values
	err = r.db.QueryRow(query,
		account.CustomerID,
		account.AccountNumber,
		account.Currency,
		account.Balance,
		account.Status,
		account.CreatedAt,
		account.UpdatedAt,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	var new_card card.CreateCardRequest
	new_card.AccountID = account.ID
	new_card.Brand = "Master Card"

	cardRepo := card.NewRepository(r.db)

	cardRepo.Create(&new_card)

	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	return account, nil
}

func (r *Repository) GetByID(id string) (*Account, error) {
	query := `
		SELECT id, customer_id, account_number, currency, balance, status, created_at, updated_at 
		FROM accounts WHERE id = $1
	`
	var a Account
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.CustomerID, &a.AccountNumber, &a.Currency, &a.Balance, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) GetByCustomerID(customerID string) ([]*Account, error) {
	query := `
		SELECT id, customer_id, account_number, currency, balance, status, created_at, updated_at 
		FROM accounts WHERE customer_id = $1
	`
	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account
	for rows.Next() {
		var a Account
		err := rows.Scan(&a.ID, &a.CustomerID, &a.AccountNumber, &a.Currency, &a.Balance, &a.Status, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &a)
	}
	return accounts, nil
}

func (r *Repository) Update(id string, a *Account) error {
	query := `
		UPDATE accounts SET balance = $1, currency = $2, status = $3, updated_at = $4 
		WHERE id = $5
	`
	result, err := r.db.Exec(query, a.Balance, a.Currency, a.Status, a.UpdatedAt, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("account not found")
	}
	return nil
}

func (r *Repository) Close(id string) error {
	query := `UPDATE accounts SET status = 'closed', updated_at = NOW() WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("account not found")
	}
	return nil
}

func (r *Repository) UpdateBalance(id string, newBalance float64) error {
	query := `UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, newBalance, id)
	return err
}

func (r *Repository) GetAll() ([]*Account, error) {
	query := `SELECT * FROM account`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account

	for rows.Next() {
		var acc Account

		err := rows.Scan(
			&acc.ID,
			&acc.AccountNumber,
			&acc.Balance,
			&acc.CreatedAt,
			&acc.Currency,
			&acc.CustomerID,
			&acc.UpdatedAt,
			&acc.Status,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, &acc)
	}

	return accounts, nil
}

func (r *Repository) Delete(id string) error {
	query := `DELETE FROM accounts WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account with id %s not found", id)
	}

	return nil
}
