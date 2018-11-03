package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Customer models a customer in the database
type Customer struct {
	ID        int        `db:"id" json:"id"`
	Msisdn    string     `db:"msisdn" json:"msisdn"`
	Ticket    string     `db:"ticket" json:"ticket"`
	QueueID   int        `db:"queue_id" json:"queueId"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	ServedBy  int        `db:"served_by" json:"servedBy"`
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
func (repo *CustomersRepo) Create(c *Customer) error {
	query := "INSERT INTO customers (msisdn, ticket, queue_id) VALUES (?, ?, ?)"

	_, err := repo.db.Exec(query, c.Msisdn, c.Ticket, c.QueueID)

	return err
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
	query := "SELECT c.* FROM customers AS c WHERE c.served_at = NULL"

	var customers []*Customer
	err := repo.db.Select(&customers, query)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

// MarkAsServed marks a customer as served in the database
func (repo *CustomersRepo) MarkAsServed(custId, userId int) error {
	query := "UPDATE customers SET served_at = CURRENT_TIMESTAMP(), served_by = ? WHERE id = ?"

	_, err := repo.db.Exec(query, custId, userId)

	return err
}
