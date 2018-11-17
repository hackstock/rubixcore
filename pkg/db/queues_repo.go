package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Queue models a queue in the db
type Queue struct {
	ID          int64      `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	IsActive    bool       `db:"is_active" json:"isActive"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
}

// QueuesRepo defines methods for executing business rules
// on queues
type QueuesRepo struct {
	db *sqlx.DB
}

// NewQueuesRepo returns a pointer to a QueuesRepo
func NewQueuesRepo(db *sqlx.DB) *QueuesRepo {
	return &QueuesRepo{db}
}

// Create saves a queue into the database
func (repo *QueuesRepo) Create(q *Queue) (*Queue, error) {
	query := "INSERT INTO queues (name, description) VALUES (?, ?)"
	res, err := repo.db.Exec(query, q.Name, q.Description)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	q.ID = id
	return q, nil
}

// GetAll fetches and returns all queues from the database
func (repo *QueuesRepo) GetAll() ([]*Queue, error) {
	query := "SELECT q.* FROM queues AS q"

	var queues []*Queue
	err := repo.db.Select(&queues, query)
	if err != nil {
		return nil, err
	}

	return queues, nil
}

// GetActive fetches and returns all active queues from the database
func (repo *QueuesRepo) GetActive() ([]*Queue, error) {
	query := "SELECT q.* FROM queues AS q WHERE q.is_active = TRUE"

	var queues []*Queue
	err := repo.db.Select(&queues, query)
	if err != nil {
		return nil, err
	}

	return queues, nil
}

// Get fetches and returns a queue by id
func (repo *QueuesRepo) Get(id int64) (*Queue, error) {
	query := "SELECT q.* FROM queues AS q WHERE q.id = ?"

	q := new(Queue)
	err := repo.db.QueryRowx(query, id).StructScan(q)
	if err != nil {
		return nil, err
	}

	return q, nil
}

// Update updates the name, descrition, or activity status
// of a queue and returns the updated record
func (repo *QueuesRepo) Update(q *Queue) (*Queue, error) {
	query := "UPDATE queues SET name = ?, description = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP() WHERE id = ?"

	_, err := repo.db.Exec(query, q.Name, q.Description, q.IsActive, q.ID)
	if err != nil {
		return nil, err
	}

	return repo.Get(q.ID)
}

// Delete removes a queue from the database by id
func (repo *QueuesRepo) Delete(id int64) error {
	query := "DELETE FROM queues WHERE id = ?"

	_, err := repo.db.Exec(query, id)

	return err
}
