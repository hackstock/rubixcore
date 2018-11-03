package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestHashPassword(t *testing.T) {
	testCases := []struct {
		tag           string
		plainPassword string
	}{
		{
			tag:           "valid case",
			plainPassword: "foobar",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tag, func(t *testing.T) {
			_, err := hashPassword(tc.plainPassword)
			if err != nil {
				t.Fatalf("expected no error hashing password, got %v", err)
			}
		})
	}
}

func TestComparePasswords(t *testing.T) {
	testCases := []struct {
		tag            string
		hashedPassword string
		plainPassword  string
		expect         bool
	}{
		{
			tag:            "valid case",
			hashedPassword: "$2a$10$pJofeBaFtdXo4RdRrKBJF.FW/ePvnS3.xgNpdC0N4FNt2S1H3QO2K",
			plainPassword:  "foobar",
			expect:         true,
		},
		{
			tag:            "invalid case",
			hashedPassword: "$2a$10$pJofeBaFtdXo4RdRrKBJF.FW/ePvnS3.xgNpdC0N4FNt2S1H3QO2K",
			plainPassword:  "somepassword",
			expect:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tag, func(t *testing.T) {
			got := comparePasswords(tc.hashedPassword, tc.plainPassword)
			if got != tc.expect {
				t.Fatalf("expected %v, got %v", tc.expect, got)
			}
		})
	}
}

func TestCreate_ShouldPass(t *testing.T) {
	query := `^INSERT INTO user_accounts \(username, password, is_admin\) VALUES \(\?, \?, \?\)$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u := &UserAccount{Username: "someuser", Password: "somepassword", IsAdmin: false}

	mock.ExpectExec(query).
		WithArgs(
			u.Username,
			u.Password,
			u.IsAdmin,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewUsersRepo(dbMock)

	err = repo.Create(u)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreate_ShouldFail(t *testing.T) {
	query := `^INSERT INTO user_accounts \(username, password, is_admin\) VALUES \(\?, \?, \?\)$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u := &UserAccount{Username: "someuser", Password: "somepassword", IsAdmin: false}

	mock.ExpectExec(query).
		WithArgs(
			u.Username,
			u.Password,
			u.IsAdmin,
		).
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewUsersRepo(dbMock)

	err = repo.Create(u)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAll_ShouldPass(t *testing.T) {
	query := `^SELECT \* FROM user_accounts$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(query).WithArgs().WillReturnRows(
		sqlmock.NewRows([]string{"id", "username", "password", "is_admin", "created_at", "last_login_at", "updated_at"}).
			AddRow(1, "user1", "$2a$10$pJofeBaFtdXo4RdRrKBJF.FW/ePvnS3.xgNpdC0N4FNt2S1H3QO2K", false, time.Now(), nil, nil).
			AddRow(2, "user2", "$2a$10$pJofeBaFtdXo4RdRrKBJF.FW/ePvnS3.xgNpdC0N4FNt2S1H3QO2K", false, time.Now(), nil, nil).
			AddRow(3, "user3", "$2a$10$pJofeBaFtdXo4RdRrKBJF.FW/ePvnS3.xgNpdC0N4FNt2S1H3QO2K", false, time.Now(), nil, nil).
			AddRow(4, "user4", "$2a$10$pJofeBaFtdXo4RdRrKBJF.FW/ePvnS3.xgNpdC0N4FNt2S1H3QO2K", false, time.Now(), nil, nil),
	)

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewUsersRepo(dbMock)

	_, err = repo.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAll_ShouldFail(t *testing.T) {
	query := `^SELECT \* FROM user_accounts$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(query).
		WithArgs().
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewUsersRepo(dbMock)

	_, err = repo.GetAll()
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthenticate_ShouldPass(t *testing.T) {
	query := `^SELECT \* FROM user_accounts AS u WHERE u.username = \? AND u.password = \?$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	username := "someuser"
	password := "somepassword"

	mock.ExpectQuery(query).WithArgs(
		username,
		password,
	).WillReturnRows(
		sqlmock.NewRows([]string{"id", "username", "password", "is_admin", "created_at", "last_login_at", "updated_at"}).
			AddRow(1, username, "$2a$10$pJofeBaFtdXo4RdRrKBJF.FW/ePvnS3.xgNpdC0N4FNt2S1H3QO2K", false, time.Now(), nil, nil),
	)

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewUsersRepo(dbMock)

	_, err = repo.Authenticate(username, password)
	if err != nil {
		t.Fatalf("expected not error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthenticate_ShouldFail(t *testing.T) {
	query := `^SELECT \* FROM user_accounts AS u WHERE u.username = \? AND u.password = \?$`

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	username := "someuser"
	password := "somepassword"

	mock.ExpectQuery(query).WithArgs(
		username,
		password,
	).
		WillReturnError(fmt.Errorf("db error"))

	dbMock := sqlx.NewDb(db, "sqlmock")
	repo := NewUsersRepo(dbMock)

	user, err := repo.Authenticate(username, password)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if user != nil {
		t.Fatalf("expected nil, got %v", user)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateLastLogin_ShouldPass(t *testing.T) {
	query := `^UPDATE user_accounts SET last_login_at = CURRENT_TIMESTAMP\(\) WHERE id = \?$`

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
	repo := NewUsersRepo(dbMock)

	err = repo.UpdateLastLogin(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateLastLogin_ShouldFail(t *testing.T) {
	query := `^UPDATE user_accounts SET last_login_at = CURRENT_TIMESTAMP\(\) WHERE id = \?$`

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
	repo := NewUsersRepo(dbMock)

	err = repo.UpdateLastLogin(id)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
