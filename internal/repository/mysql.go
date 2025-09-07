package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/domain"
)

type MySqlPersonRepo struct {
	db *sql.DB
}

func NewMySqlPersonRepo(db *sql.DB) *MySqlPersonRepo {
	return &MySqlPersonRepo{db: db}
}

// CreatePerson inserts a new person
func (r *MySqlPersonRepo) CreatePerson(ctx context.Context, person domain.Person) error {
	jsonHobbies, err := json.Marshal(person.Hobbies)
	if err != nil {
		return err
	}

	query := `INSERT INTO person (id, name, age, hobbies) VALUES (?, ?, ?, ?)`
	_, err = r.db.ExecContext(ctx, query, person.Id, person.Name, person.Age, jsonHobbies)
	if err != nil {
		return err
	}

	return nil
}

// GetAllPersons retrieves persons with pagination
func (r *MySqlPersonRepo) GetAllPersons(ctx context.Context, limit, offset int) ([]domain.Person, error) {
	query := `SELECT id, name, age, hobbies FROM person LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []domain.Person

	for rows.Next() {
		var p domain.Person
		var hobbiesJSON []byte

		if err := rows.Scan(&p.Id, &p.Name, &p.Age, &hobbiesJSON); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(hobbiesJSON, &p.Hobbies); err != nil {
			return nil, err
		}

		persons = append(persons, p)
	}

	// check if rows iteration itself caused error
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

// GetPerson retrieves a single person by ID
func (r *MySqlPersonRepo) GetPerson(ctx context.Context, id uuid.UUID) (*domain.Person, error) {
	query := `SELECT id, name, age, hobbies FROM person WHERE id=?`

	row := r.db.QueryRowContext(ctx, query, id)

	var p domain.Person
	var hobbiesJSON []byte

	if err := row.Scan(&p.Id, &p.Name, &p.Age, &hobbiesJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("person not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(hobbiesJSON, &p.Hobbies); err != nil {
		return nil, err
	}

	return &p, nil
}

// UpdatePerson updates an existing person
func (r *MySqlPersonRepo) UpdatePerson(ctx context.Context, person domain.Person) error {
	hobbiesJSON, err := json.Marshal(person.Hobbies)
	if err != nil {
		return err
	}

	query := `UPDATE person SET name=?, age=?, hobbies=? WHERE id=?`
	res, err := r.db.ExecContext(ctx, query, person.Name, person.Age, hobbiesJSON, person.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("person not found")
	}

	return nil
}

// DeletePerson deletes a person by ID
func (r *MySqlPersonRepo) DeletePerson(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM person WHERE id=?`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("person not found")
	}

	return nil
}
