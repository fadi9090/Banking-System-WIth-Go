package person

import (
	"database/sql"
	"task-project/internal/utility"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(p *CreatePersonRequest) (*Person, error) {
	const query = `
        INSERT INTO persons (first_name, last_name, email, person_number, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING id, first_name, last_name, email, person_number, created_at, updated_at
    `

	var person Person

	// Use QueryRow instead of Exec when you need to RETURN values
	err := r.db.QueryRow(
		query,
		p.FirstName,
		p.LastName,
		p.Email,
		utility.GeneratePersonNumber(),
		time.Now(),
		time.Now(),
	).Scan(
		&person.ID,
		&person.FirstName,
		&person.LastName,
		&person.Email,
		&person.PersonNumber,
		&person.Created_at,
		&person.Updated_at,
	)

	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (r *Repository) Update(id string, p *UpdatePersonRequest) error {
	const query = "UPDATE persons SET first_name = $1, last_name = $2, email = $3, updated_at = $4 WHERE id = $5"

	result, err := r.db.Exec(query, p.FirstName, p.LastName, p.Email, time.Now(), id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) Delete(id string) error {
	const query = "DELETE FROM persons WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) GetByID(id string) (*Person, error) {
	const query = "SELECT id, first_name, last_name, email, person_number, created_at, updated_at FROM persons WHERE id = $1"

	var person Person
	err := r.db.QueryRow(query, id).Scan(
		&person.ID,
		&person.FirstName,
		&person.LastName,
		&person.Email,
		&person.PersonNumber,
		&person.Created_at,
		&person.Updated_at,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *Repository) GetByEmail(email string) (*Person, error) {
	const query = "SELECT id, first_name, last_name, email, person_number, created_at, updated_at FROM persons WHERE email = $1"

	var person Person
	err := r.db.QueryRow(query, email).Scan(
		&person.ID,
		&person.FirstName,
		&person.LastName,
		&person.Email,
		&person.PersonNumber,
		&person.Created_at,
		&person.Updated_at,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *Repository) GetByNumber(num string) (*Person, error) {
	const query = "SELECT id, first_name, last_name, email, person_number, created_at, updated_at FROM persons WHERE person_number = $1"

	var person Person
	err := r.db.QueryRow(query, num).Scan(
		&person.ID,
		&person.FirstName,
		&person.LastName,
		&person.Email,
		&person.PersonNumber,
		&person.Created_at,
		&person.Updated_at,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &person, nil
}
