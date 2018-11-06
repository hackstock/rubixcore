package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestCreateQueue_ShouldPass(t *testing.T) {
	query := `^INSERT INTO queues \(name, description\) VALUES \(\?, \?\)$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	q := &Queue{Name: "testqueue", Description: "test queue description"}

	mock.ExpectExec(query).
		WithArgs(
			q.Name,
			q.Description,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock := sqlx.NewDb(db, "sqlmock")
	queuesRepo := NewQueuesRepo(dbMock)

	err = queuesRepo.Create(q)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateQueue_ShouldFail(t *testing.T) {
	query := `^INSERT INTO queues \(name, description\) VALUES \(\?, \?\)$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	q := &Queue{Name: "testqueue", Description: "test queue description"}

	mock.ExpectExec(query).
		WithArgs(
			q.Name,
			q.Description,
		).
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	queuesRepo := NewQueuesRepo(dbMock)

	err = queuesRepo.Create(q)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllQueue_ShouldPass(t *testing.T) {
	query := `^SELECT q.\* FROM queues AS q$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(query).WithArgs().WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "description", "is_active", "created_at"}).
			AddRow(1, "Queue name", "Queue descrition", false, time.Now()).
			AddRow(2, "Queue name", "Queue descrition", true, time.Now()).
			AddRow(3, "Queue name", "Queue descrition", false, time.Now()).
			AddRow(4, "Queue name", "Queue descrition", true, time.Now()),
	)

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewQueuesRepo(dbMock)

	_, err = repo.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllQueue_ShouldFail(t *testing.T) {
	query := `^SELECT q.\* FROM queues AS q$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(query).
		WithArgs().
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	queuesRepo := NewQueuesRepo(dbMock)

	queues, err := queuesRepo.GetAll()
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if queues != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllActiveQueues_ShouldPass(t *testing.T) {
	query := `^SELECT q.\* FROM queues AS q WHERE q.is_active = TRUE$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(query).WithArgs().WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "description", "is_active", "created_at"}).
			AddRow(1, "Queue name", "Queue descrition", true, time.Now()).
			AddRow(2, "Queue name", "Queue descrition", true, time.Now()).
			AddRow(3, "Queue name", "Queue descrition", true, time.Now()).
			AddRow(4, "Queue name", "Queue descrition", true, time.Now()),
	)

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewQueuesRepo(dbMock)

	_, err = repo.GetActive()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllActiveQueues_ShouldFail(t *testing.T) {
	query := `^SELECT q.\* FROM queues AS q WHERE q.is_active = TRUE$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(query).
		WithArgs().
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	queuesRepo := NewQueuesRepo(dbMock)

	queues, err := queuesRepo.GetActive()
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if queues != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetQueue_ShouldPass(t *testing.T) {
	query := `^SELECT q.\* FROM queues AS q WHERE q.id = \?$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	id := 1

	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "description", "is_active", "created_at"}).
			AddRow(1, "Queue name", "Queue descrition", false, time.Now()),
	)

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewQueuesRepo(dbMock)

	_, err = repo.Get(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetQueue_ShouldFail(t *testing.T) {
	query := `^SELECT q.\* FROM queues AS q WHERE q.id = \?$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := 1

	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	queuesRepo := NewQueuesRepo(dbMock)

	queue, err := queuesRepo.Get(id)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if queue != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteQueue_ShouldPass(t *testing.T) {
	query := `^DELETE FROM queues WHERE id = \?$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := 1

	mock.ExpectExec(query).
		WithArgs(
			id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock := sqlx.NewDb(db, "sqlmock")
	queuesRepo := NewQueuesRepo(dbMock)

	err = queuesRepo.Delete(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteQueue_ShouldFail(t *testing.T) {
	query := `^DELETE FROM queues WHERE id = \?$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := 1

	mock.ExpectExec(query).
		WithArgs(
			id,
		).
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	queuesRepo := NewQueuesRepo(dbMock)

	err = queuesRepo.Delete(id)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateQueue_ShouldPass(t *testing.T) {
	query := `^UPDATE queues SET name = \?, description = \?, is_active = \?, updated_at = CURRENT_TIMESTAMP\(\) WHERE id = \?$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	q := &Queue{ID: 1, Name: "test queue", Description: "test description", IsActive: true}

	mock.ExpectExec(query).
		WithArgs(
			q.Name,
			q.Description,
			q.IsActive,
			q.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(`^SELECT q.\* FROM queues AS q WHERE q.id = \?$`).WithArgs(q.ID).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "description", "is_active"}).
			AddRow(q.ID, q.Name, q.Description, q.IsActive),
	)

	dbMock := sqlx.NewDb(db, "sqlmock")
	queuesRepo := NewQueuesRepo(dbMock)

	updatedQueue, err := queuesRepo.Update(q)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedQueue == nil {
		t.Fatalf("expected %v, got nil", q)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateQueue_ShouldFail(t *testing.T) {
	query := `^UPDATE queues SET name = \?, description = \?, is_active = \?, updated_at = CURRENT_TIMESTAMP\(\) WHERE id = \?$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	q := &Queue{ID: 1, Name: "test queue", Description: "test description", IsActive: true}

	mock.ExpectExec(query).
		WithArgs(
			q.Name,
			q.Description,
			q.IsActive,
			q.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(`^SELECT q.\* FROM queues AS q WHERE q.id = \?$`).WithArgs(q.ID).
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	queuesRepo := NewQueuesRepo(dbMock)

	updatedQueue, err := queuesRepo.Update(q)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if updatedQueue != nil {
		t.Fatalf("expected nil, got %v", updatedQueue)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
