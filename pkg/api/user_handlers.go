package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hackstock/rubixcore/pkg/db"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func getAllUsers(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := db.NewUsersRepo(dbConn)
		users, err := repo.GetAll()
		if err != nil {
			handleServerError(w, "failed fetching all users", err, logger)
			return
		}

		render.JSON(w, r, Response{Data: users})
	}
}

func createUser(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := new(db.UserAccount)
		err := json.NewDecoder(r.Body).Decode(account)
		if err != nil {
			handleBadRequest(w, "failed decoding request payload", err, logger)
			return
		}

		hash, err := hashPassword(account.Password)
		if err != nil {
			handleServerError(w, "failed hasing account password", err, logger)
			return
		}
		account.Password = hash

		repo := db.NewUsersRepo(dbConn)
		err = repo.Create(account)
		if err != nil {
			handleServerError(w, "failed saving user into db", err, logger)
			return
		}

		render.JSON(w, r, Response{Info: "account created successfully"})
	}
}

func authenticate(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials = struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			handleServerError(w, "failed decoding request payload", err, logger)
			return
		}

		repo := db.NewUsersRepo(dbConn)
		account, err := repo.GetByUsername(credentials.Username)
		logger.Info("found user", zap.Any("user", account))
		if err != nil {
			handleServerError(w, "failed fetching account", err, logger)
			return
		}

		if comparePasswords(account.Password, credentials.Password) == false {
			handleServerError(w, "failed comparing passwords", err, logger)
			return
		}

		render.JSON(w, r, Response{Data: account, Info: "authentication succeeded"})
	}
}

func usersRoutes(dbConn *sqlx.DB, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", getAllUsers(dbConn, logger))
	router.Post("/", createUser(dbConn, logger))
	router.Post("/login", authenticate(dbConn, logger))

	return router
}
