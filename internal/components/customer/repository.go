package customer

import (
	"database/sql"
	"errors"
	"task-project/internal/utility"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(c *CreateCustomerRequest) (*Customer, error) {
	query := `
		INSERT INTO customers (username, email, person_id, customer_number, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, email, person_id, customer_number, created_at, updated_at
	`

	customer := &Customer{}
	err := r.db.QueryRow(
		query,
		c.Username,
		c.Email,
		c.PersonID,
		utility.GenerateCustomerNumber(),
		time.Now(),
		time.Now(),
	).Scan(
		&customer.ID,
		&customer.Username,
		&customer.Email,
		&customer.PersonID,
		&customer.CustomerNumber,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (r *Repository) GetByID(id string) (*Customer, error) {
	query := `
		SELECT id, username, email, person_id, customer_number, created_at, updated_at 
		FROM customers WHERE id = $1
	`
	var c Customer
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.Username, &c.Email, &c.PersonID, &c.CustomerNumber, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) GetByCustomerNumber(number string) (*Customer, error) {
	query := `
		SELECT id, username, email, person_id, customer_number, created_at, updated_at 
		FROM customers WHERE customer_number = $1
	`
	var c Customer
	err := r.db.QueryRow(query, number).Scan(
		&c.ID, &c.Username, &c.Email, &c.PersonID, &c.CustomerNumber, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) Update(id string, c *Customer) error {
	query := `
		UPDATE customers SET username = $1, email = $2, updated_at = $3
		WHERE id = $4
	`
	result, err := r.db.Exec(query, c.Username, c.Email, c.UpdatedAt, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("customer not found")
	}
	return nil
}

func (r *Repository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM customers WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("customer not found")
	}
	return nil
}

func (r *Repository) GetAll() ([]*Customer, error) {
	query := `SELECT * FROM Customers`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*Customer

	for rows.Next() {
		var cus Customer

		err := rows.Scan(
			&cus.ID,
			&cus.CustomerNumber,
			&cus.Email,
			&cus.PersonID,
			&cus.CreatedAt,
			&cus.Username,
			&cus.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		customers = append(customers, &cus)
	}

	return customers, nil

}

func (r *Repository) UniqueEmail(email string) (bool, error) {
	query := `SELECT email FROM customers WHERE email = $1`

	var found string
	err := r.db.QueryRow(query, email).Scan(&found)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

func (r *Repository) UniqueUsername(username string) (bool, error) {
	query := `SELECT username FROM customers WHERE username = $1`

	var found string
	err := r.db.QueryRow(query, username).Scan(&found)

	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	return false, nil
}
