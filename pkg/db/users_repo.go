package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// UserAccount models a user in the db
type UserAccount struct {
	ID          int        `db:"id" json:"id"`
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

// Create saves a UserAccount into the database
func (repo *UsersRepo) Create(u *UserAccount) error {
	query := "INSERT INTO user_accounts (username, password, is_admin) VALUES (?, ?, ?)"

	_, err := repo.db.Exec(query, u.Username, u.Password, u.IsAdmin)

	return err
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

func hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)) == nil
}
