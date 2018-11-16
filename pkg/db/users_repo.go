package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// UserAccount models a user in the db
type UserAccount struct {
	ID          int64      `db:"id" json:"id"`
	Username    string     `db:"username" json:"username"`
	Password    string     `db:"password" json:"-"`
	IsAdmin     bool       `db:"is_admin" json:"isAdmin"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
	LastLoginAt *time.Time `db:"last_login_at" json:"lastLoginAt"`
}

// UsersRepo defines methods for running business rules on user accounts
type UsersRepo struct {
	db *sqlx.DB
}

// NewUsersRepo returns a pointer to a UsersRepo
func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{db}
}

// GetByUsername returns a user with the specified username
func (repo *UsersRepo) GetByUsername(username string) (*UserAccount, error) {
	query := "SELECT * FROM user_accounts AS u WHERE u.username = ?"

	u := new(UserAccount)
	err := repo.db.QueryRowx(query, username).StructScan(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// UpdateLastLogin updates the last login date for the specified user
func (repo *UsersRepo) UpdateLastLogin(id int64) error {
	query := "UPDATE user_accounts SET last_login_at = CURRENT_TIMESTAMP() WHERE id = ?"

	_, err := repo.db.Exec(query, id)

	return err
}

// Create saves a UserAccount into the database
func (repo *UsersRepo) Create(u *UserAccount) (*UserAccount, error) {
	query := "INSERT INTO user_accounts (username, password, is_admin) VALUES (?, ?, ?)"

	res, err := repo.db.Exec(query, u.Username, u.Password, u.IsAdmin)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

// GetAll fetches and returns all user accounts in the database
func (repo *UsersRepo) GetAll() ([]*UserAccount, error) {
	var accounts []*UserAccount
	query := "SELECT * FROM user_accounts"

	err := repo.db.Select(&accounts, query)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
