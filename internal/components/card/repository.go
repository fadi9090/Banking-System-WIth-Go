package card

import (
	"database/sql"
	"errors"
	"fmt"
	"task-project/internal/utility"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(c *CreateCardRequest) (*Card, error) {
	query := `
        INSERT INTO cards (id, account_id, last4, brand, is_blocked, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `
	last4 := utility.GenerateCardLast4()

	card := &Card{
		AccountID: c.AccountID,
		Last4:     last4,
		Brand:     c.Brand,
		IsBlocked: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := r.db.QueryRow(query,
		card.AccountID,
		card.Last4,
		card.Brand,
		card.IsBlocked,
		card.CreatedAt,
		card.UpdatedAt,
	).Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error creating card: %w", err)
	}

	return card, nil
}

func (r *Repository) GetByID(id string) (*Card, error) {
	query := `
		SELECT id, account_id, last4, brand, is_blocked, created_at, updated_at 
		FROM cards WHERE id = $1
	`
	var c Card
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.AccountID, &c.Last4, &c.Brand, &c.IsBlocked, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) GetByAccountID(accountID string) ([]*Card, error) {
	query := `
		SELECT id, account_id, last4, brand, is_blocked, created_at, updated_at 
		FROM cards WHERE account_id = $1
	`
	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*Card
	for rows.Next() {
		var c Card
		err := rows.Scan(&c.ID, &c.AccountID, &c.Last4, &c.Brand, &c.IsBlocked, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cards = append(cards, &c)
	}
	return cards, nil
}

func (r *Repository) Update(id string, c *Card) error {
	query := `
		UPDATE cards SET brand = $1, is_blocked = $2, updated_at = $3 
		WHERE id = $4
	`
	result, err := r.db.Exec(query, c.Brand, c.IsBlocked, c.UpdatedAt, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("card not found")
	}
	return nil
}

func (r *Repository) Block(id string, blocked bool) error {
	query := `UPDATE cards SET is_blocked = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(query, blocked, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("card not found")
	}
	return nil
}
