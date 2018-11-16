package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Customer models a customer in the database
type Customer struct {
	ID        int64      `db:"id" json:"id"`
	Msisdn    string     `db:"msisdn" json:"msisdn"`
	Ticket    string     `db:"ticket" json:"ticket"`
	QueueID   int64      `db:"queue_id" json:"queueId"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	ServedAt  *time.Time `db:"served_at" json:"servedAt"`
}

// CustomersRepo defines methods for executing business rules
// on customers in the database
type CustomersRepo struct {
	db *sqlx.DB
}

// NewCustomersRepo returns a pointer to a CustomersRepo
func NewCustomersRepo(db *sqlx.DB) *CustomersRepo {
	return &CustomersRepo{db}
}

// Create saves a customer into the database
func (repo *CustomersRepo) Create(c *Customer) (*Customer, error) {
	query := "INSERT INTO customers (msisdn, ticket, queue_id) VALUES (?, ?, ?)"

	res, err := repo.db.Exec(query, c.Msisdn, c.Ticket, c.QueueID)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	c.ID = id
	return c, nil
}

// GetAll fetches and return all customers from the database
func (repo *CustomersRepo) GetAll() ([]*Customer, error) {
	query := "SELECT c.* FROM customers AS c"

	var customers []*Customer
	err := repo.db.Select(&customers, query)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

// GetUnserved fetches and return all customers from the database
// that have not been served yet
func (repo *CustomersRepo) GetUnserved() ([]*Customer, error) {
	query := "SELECT c.* FROM customers AS c WHERE c.served_at IS NULL"

	var customers []*Customer
	err := repo.db.Select(&customers, query)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

// MarkAsServed marks a customer as served in the database
func (repo *CustomersRepo) MarkAsServed(custID int) error {
	query := "UPDATE customers SET served_at = NOW() WHERE id = ?"

	_, err := repo.db.Exec(query, custID)

	return err
}
