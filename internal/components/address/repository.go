package address

import (
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(a *Address) error {
	query := `
		INSERT INTO addresses (street, city, state, zip_code, country, customer_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, a.Street, a.City, a.State, a.ZipCode, a.Country, a.CustomerID, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *Repository) GetByID(id string) (*Address, error) {
	query := `
		SELECT id, street, city, state, zip_code, country, customer_id, created_at, updated_at 
		FROM addresses WHERE id = $1
	`
	var a Address
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.Street, &a.City, &a.State, &a.ZipCode, &a.Country, &a.CustomerID, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) GetByCustomerID(customerID string) ([]*Address, error) {
	query := `
		SELECT id, street, city, state, zip_code, country, customer_id, created_at, updated_at 
		FROM addresses WHERE customer_id = $1
	`
	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*Address
	for rows.Next() {
		var a Address
		err := rows.Scan(&a.ID, &a.Street, &a.City, &a.State, &a.ZipCode, &a.Country, &a.CustomerID, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, &a)
	}
	return addresses, nil
}

func (r *Repository) Update(id string, a *Address) error {
	query := `
		UPDATE addresses SET street = $1, city = $2, state = $3, zip_code = $4, country = $5, updated_at = $6 
		WHERE id = $7
	`
	result, err := r.db.Exec(query, a.Street, a.City, a.State, a.ZipCode, a.Country, a.UpdatedAt, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("address not found")
	}
	return nil
}

func (r *Repository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM addresses WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("address not found")
	}
	return nil
}
