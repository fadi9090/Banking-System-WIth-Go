package transaction

import (
	"database/sql"
	"errors"
	"task-project/internal/components/account"
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

func (r *Repository) CreateTransaction(tx *CreateTransactionRequest) error {
	query := `INSERT INTO transactions (account_id, card_id, direction, amount, description, status, currency, reference, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	// Account and balance validation
	var acc account.Account
	err := r.db.QueryRow("SELECT id, balance FROM accounts WHERE id = $1", tx.AccountID).Scan(&acc.ID, &acc.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("account not found")
		}
		return err
	}

	// Card validation (if card_id provided)
	if tx.CardID != nil && *tx.CardID != "" {
		var card card.Card
		err = r.db.QueryRow("SELECT id FROM cards WHERE id = $1 AND account_id = $2", tx.CardID, tx.AccountID).Scan(&card.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("card not found or does not belong to this account")
			}
			return err
		}
	}

	// Validate transaction (sufficient balance for debit)
	validationMsg, valid := utility.ValidateTransaction(acc.Balance, tx.Amount, tx.Direction)
	if !valid {
		return errors.New(validationMsg)
	}

	// Calculate new balance
	newBalance := acc.Balance
	if tx.Direction == "credit" {
		newBalance += tx.Amount
	} else {
		newBalance -= tx.Amount
	}

	// Start database transaction for atomic operation
	dbTx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Insert transaction
	_, err = dbTx.Exec(query, tx.AccountID, tx.CardID, tx.Direction, tx.Amount, tx.Description, "completed", tx.Currency, utility.GenerateReferenseNumber(), time.Now())
	if err != nil {
		dbTx.Rollback()
		return err
	}

	// Update account balance
	_, err = dbTx.Exec("UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2", newBalance, tx.AccountID)
	if err != nil {
		dbTx.Rollback()
		return err
	}

	return dbTx.Commit()
}

func (r *Repository) GetTranById(id string) (*Transaction, error) {
	query := `SELECT id, account_id, amount, card_id, created_at, description, status, direction, currency, reference 
			  FROM transactions WHERE id = $1`

	var tran Transaction
	err := r.db.QueryRow(query, id).Scan(
		&tran.ID,
		&tran.AccountID,
		&tran.Amount,
		&tran.CardID,
		&tran.CreatedAt,
		&tran.Description,
		&tran.Status,
		&tran.Direction,
		&tran.Currency,
		&tran.Reference,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &tran, nil
}

func (r *Repository) GetTranByAccId(id string) ([]*Transaction, error) {
	query := `SELECT id, account_id, amount, card_id, created_at, description, status, direction, currency, reference 
			  FROM transactions WHERE account_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var tran Transaction

		err := rows.Scan(
			&tran.ID,
			&tran.AccountID,
			&tran.Amount,
			&tran.CardID,
			&tran.CreatedAt,
			&tran.Description,
			&tran.Status,
			&tran.Direction,
			&tran.Currency,
			&tran.Reference,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &tran)
	}

	return transactions, nil
}

func (r *Repository) GetTranByAmount(amount float64) ([]*Transaction, error) {
	query := `SELECT id, account_id, amount, card_id, created_at, description, status, direction, currency, reference 
			  FROM transactions WHERE amount = $1`

	rows, err := r.db.Query(query, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var tran Transaction

		err := rows.Scan(
			&tran.ID,
			&tran.AccountID,
			&tran.Amount,
			&tran.CardID,
			&tran.CreatedAt,
			&tran.Description,
			&tran.Status,
			&tran.Direction,
			&tran.Currency,
			&tran.Reference,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &tran)
	}

	return transactions, nil
}

func (r *Repository) GetTranByCardId(id string) ([]*Transaction, error) {
	query := `SELECT id, account_id, amount, card_id, created_at, description, status, direction, currency, reference 
			  FROM transactions WHERE card_id = $1`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var tran Transaction

		err := rows.Scan(
			&tran.ID,
			&tran.AccountID,
			&tran.Amount,
			&tran.CardID,
			&tran.CreatedAt,
			&tran.Description,
			&tran.Status,
			&tran.Direction,
			&tran.Currency,
			&tran.Reference,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &tran)
	}

	return transactions, nil
}

func (r *Repository) GetTranByStatus(status string) ([]*Transaction, error) {
	query := `SELECT id, account_id, amount, card_id, created_at, description, status, direction, currency, reference 
			  FROM transactions WHERE status = $1`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var tran Transaction

		err := rows.Scan(
			&tran.ID,
			&tran.AccountID,
			&tran.Amount,
			&tran.CardID,
			&tran.CreatedAt,
			&tran.Description,
			&tran.Status,
			&tran.Direction,
			&tran.Currency,
			&tran.Reference,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &tran)
	}

	return transactions, nil
}

func (r *Repository) GetTranByDirection(dir string) ([]*Transaction, error) {
	query := `SELECT id, account_id, amount, card_id, created_at, description, status, direction, currency, reference 
			  FROM transactions WHERE direction = $1`

	rows, err := r.db.Query(query, dir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var tran Transaction

		err := rows.Scan(
			&tran.ID,
			&tran.AccountID,
			&tran.Amount,
			&tran.CardID,
			&tran.CreatedAt,
			&tran.Description,
			&tran.Status,
			&tran.Direction,
			&tran.Currency,
			&tran.Reference,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &tran)
	}

	return transactions, nil
}

func (r *Repository) GetTranByDate(date time.Time) ([]*Transaction, error) {
	query := `SELECT id, account_id, amount, card_id, created_at, description, status, direction, currency, reference 
			  FROM transactions WHERE DATE(created_at) = DATE($1)`

	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var tran Transaction

		err := rows.Scan(
			&tran.ID,
			&tran.AccountID,
			&tran.Amount,
			&tran.CardID,
			&tran.CreatedAt,
			&tran.Description,
			&tran.Status,
			&tran.Direction,
			&tran.Currency,
			&tran.Reference,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &tran)
	}

	return transactions, nil
}

func (r *Repository) GetAllTransactions() ([]*Transaction, error) {
	query := `SELECT id, account_id, amount, card_id, created_at, description, status, direction, currency, reference 
			  FROM transactions ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var tran Transaction

		err := rows.Scan(
			&tran.ID,
			&tran.AccountID,
			&tran.Amount,
			&tran.CardID,
			&tran.CreatedAt,
			&tran.Description,
			&tran.Status,
			&tran.Direction,
			&tran.Currency,
			&tran.Reference,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &tran)
	}

	return transactions, nil
}
