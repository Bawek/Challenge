package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/domain"
)

type PostgresPersonRepo struct {
	db *pgxpool.Pool
}

func NewPostgresPersonRepo(db *pgxpool.Pool) *PostgresPersonRepo {
	return &PostgresPersonRepo{db: db}
}

// Create
func (r *PostgresPersonRepo) CreatePerson(ctx context.Context, person domain.Person) error {
	hobbiesJSON, _ := json.Marshal(person.Hobbies)

	_, err := r.db.Exec(ctx,
		`INSERT INTO persons (id, name, age, hobbies) VALUES ($1, $2, $3, $4)`,
		person.Id, person.Name, person.Age, hobbiesJSON)
	return err
}

// Get All
func (r *PostgresPersonRepo) GetAllPersons(ctx context.Context, limit, offset int) ([]domain.Person, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, age, hobbies FROM persons LIMIT $1 OFFSET $2`,
		limit, offset)
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
		_ = json.Unmarshal(hobbiesJSON, &p.Hobbies)
		persons = append(persons, p)
	}
	return persons, nil
}

// Update
func (r *PostgresPersonRepo) UpdatePerson(ctx context.Context, person domain.Person) error {
	hobbiesJSON, _ := json.Marshal(person.Hobbies)

	cmd, err := r.db.Exec(ctx,
		`UPDATE persons SET name=$1, age=$2, hobbies=$3 WHERE id=$4`,
		person.Name, person.Age, hobbiesJSON, person.Id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("person not found")
	}
	return nil
}

// Delete
func (r *PostgresPersonRepo) DeletePerson(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.db.Exec(ctx, `DELETE FROM persons WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("person not found")
	}
	return nil
}

// Get by ID
func (r *PostgresPersonRepo) GetPerson(ctx context.Context, id uuid.UUID) (*domain.Person, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, age, hobbies FROM persons WHERE id=$1`, id)

	var p domain.Person
	var hobbiesJSON []byte
	if err := row.Scan(&p.Id, &p.Name, &p.Age, &hobbiesJSON); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(hobbiesJSON, &p.Hobbies)
	return &p, nil
}
