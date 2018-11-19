package api

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
		var payload = struct {
			Username string `json:"username"`
			Password string `json:"password"`
			IsAdmin  bool   `json:"isAdmin"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			handleBadRequest(w, "failed decoding request payload", err, logger)
			return
		}

		account := &db.UserAccount{
			Username: payload.Username,
			Password: payload.Password,
			IsAdmin:  payload.IsAdmin,
		}

		hash, err := hashPassword(account.Password)
		if err != nil {
			handleServerError(w, "failed hasing account password", err, logger)
			return
		}
		account.Password = hash

		repo := db.NewUsersRepo(dbConn)
		u, err := repo.Create(account)
		if err != nil {
			handleServerError(w, "failed saving user into db", err, logger)
			return
		}

		render.JSON(w, r, Response{Data: u, Info: "account created successfully"})
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
		if err != nil {
			handleServerError(w, "failed fetching account", err, logger)
			return
		}

		if comparePasswords(account.Password, credentials.Password) == false {
			handleServerError(w, "failed comparing passwords", err, logger)
			return
		}
		token, err := generateJWT()
		if err != nil {
			handleServerError(w, "failed generating JWT", err, logger)
			return
		}

		responsePayload := struct {
			User  interface{} `json:"user"`
			Token string      `json:"token"`
		}{
			account,
			token,
		}
		render.JSON(w, r, Response{Data: responsePayload, Info: "authentication succeeded"})
	}
}

func generateJWT() (string, error) {
	now := time.Now()
	claims := jwt.StandardClaims{
		Issuer:    "rubix",
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(24 * 7 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(`secret`))
}

func usersRoutes(dbConn *sqlx.DB, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", getAllUsers(dbConn, logger))
	router.Post("/", createUser(dbConn, logger))
	router.Post("/login", authenticate(dbConn, logger))

	return router
}
