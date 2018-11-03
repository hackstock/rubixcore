package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestCreateCustomer_ShouldPass(t *testing.T) {
	query := `^INSERT INTO customers \(msisdn, ticket, queue_id\) VALUES \(\?, \?, \?\)$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	c := &Customer{Msisdn: "+233200662782", Ticket: "A201", QueueID: 1}

	mock.ExpectExec(query).
		WithArgs(
			c.Msisdn,
			c.Ticket,
			c.QueueID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock := sqlx.NewDb(db, "sqlmock")
	customersRepo := NewCustomersRepo(dbMock)

	err = customersRepo.Create(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateCustomer_ShouldFail(t *testing.T) {
	query := `^INSERT INTO customers \(msisdn, ticket, queue_id\) VALUES \(\?, \?, \?\)$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	c := &Customer{Msisdn: "+233200662782", Ticket: "A201", QueueID: 1}

	mock.ExpectExec(query).
		WithArgs(
			c.Msisdn,
			c.Ticket,
			c.QueueID,
		).
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	customersRepo := NewCustomersRepo(dbMock)

	err = customersRepo.Create(c)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllCustomer_ShouldPass(t *testing.T) {
	query := `^SELECT c.\* FROM customers AS c$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(query).WithArgs().WillReturnRows(
		sqlmock.NewRows([]string{"id", "msisdn", "ticket", "queue_id", "created_at", "served_by", "served_at"}).
			AddRow(1, "+233200662782", "A101", 1, time.Now(), 1, time.Now()).
			AddRow(2, "+233200662783", "A201", 2, time.Now(), 2, time.Now()).
			AddRow(3, "+233200662784", "A103", 3, time.Now(), 1, time.Now()).
			AddRow(4, "+233200662785", "A141", 1, time.Now(), 2, time.Now()),
	)

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewCustomersRepo(dbMock)

	customers, err := repo.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if customers == nil {
		t.Fatalf("expected list of customers, got nil")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllCustomer_ShouldFail(t *testing.T) {
	query := `^SELECT c.\* FROM customers AS c$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(query).
		WithArgs().
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewCustomersRepo(dbMock)

	customers, err := repo.GetAll()
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if customers != nil {
		t.Fatalf("expected nil , got %v", customers)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllCustomerUnserved_ShouldPass(t *testing.T) {
	query := `^SELECT c.\* FROM customers AS c WHERE c.served_at = NULL$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(query).WithArgs().WillReturnRows(
		sqlmock.NewRows([]string{"id", "msisdn", "ticket", "queue_id", "created_at", "served_by", "served_at"}).
			AddRow(1, "+233200662782", "A101", 1, time.Now(), 1, nil).
			AddRow(2, "+233200662783", "A201", 2, time.Now(), 2, nil).
			AddRow(3, "+233200662784", "A103", 3, time.Now(), 1, nil).
			AddRow(4, "+233200662785", "A141", 1, time.Now(), 2, nil),
	)

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewCustomersRepo(dbMock)

	customers, err := repo.GetUnserved()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if customers == nil {
		t.Fatalf("expected list of customers, got nil")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllCustomerUnserved_ShouldFail(t *testing.T) {
	query := `^SELECT c.\* FROM customers AS c WHERE c.served_at = NULL$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(query).
		WithArgs().
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewCustomersRepo(dbMock)

	customers, err := repo.GetUnserved()
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if customers != nil {
		t.Fatalf("expected nil , got %v", customers)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMarkAsServedCustomer_ShouldPass(t *testing.T) {
	query := `^UPDATE customers SET served_at = CURRENT_TIMESTAMP\(\), served_by = \? WHERE id = \?$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	custId := 1
	userId := 2

	mock.ExpectExec(query).
		WithArgs(
			custId,
			userId,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock := sqlx.NewDb(db, "sqlmock")
	customersRepo := NewCustomersRepo(dbMock)

	err = customersRepo.MarkAsServed(custId, userId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMarkAsServedCustomer_ShouldFail(t *testing.T) {
	query := `^UPDATE customers SET served_at = CURRENT_TIMESTAMP\(\), served_by = \? WHERE id = \?$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	custId := 1
	userId := 2

	mock.ExpectExec(query).
		WithArgs(
			custId,
			userId,
		).
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	customersRepo := NewCustomersRepo(dbMock)

	err = customersRepo.MarkAsServed(custId, userId)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
