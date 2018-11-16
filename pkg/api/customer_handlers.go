package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hackstock/rubixcore/pkg/app"
	"github.com/hackstock/rubixcore/pkg/db"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func createCustomer(rubix *app.Rubix, dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customer db.Customer
		err := json.NewDecoder(r.Body).Decode(&customer)
		if err != nil {
			handleServerError(w, "failed decoding request payload", err, logger)
			return
		}

		repo := db.NewCustomersRepo(dbConn)
		customer.Ticket = rubix.GenerateTicket()

		c, err := repo.Create(&customer)
		if err != nil {
			handleServerError(w, "failed creating customer", err, logger)
			return
		}

		rubix.AddCustomerToWaitList(customer.QueueID, customer.Msisdn, customer.Ticket)

		render.JSON(w, r, Response{Data: c, Info: "customer created successfully"})
	}
}

func getAllCustomers(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := db.NewCustomersRepo(dbConn)
		customers, err := repo.GetAll()
		if err != nil {
			handleServerError(w, "failed fetching all customers", err, logger)
			return
		}

		render.JSON(w, r, Response{Data: customers})
	}
}

func getUnservedCustomers(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := db.NewCustomersRepo(dbConn)
		customers, err := repo.GetUnserved()
		if err != nil {
			handleServerError(w, "failed fetching unserved customers", err, logger)
			return
		}

		render.JSON(w, r, Response{Data: customers})
	}
}

func markAsServed(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			CustomerID int `json:"customerId"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			handleServerError(w, "failed decoding request payload", err, logger)
			return
		}

		repo := db.NewCustomersRepo(dbConn)
		err = repo.MarkAsServed(payload.CustomerID)
		if err != nil {
			handleServerError(w, "failed marking customer as served", err, logger)
			return
		}

		render.JSON(w, r, Response{Info: "customer called successfully"})
	}
}

func customersRoutes(rubix *app.Rubix, dbConn *sqlx.DB, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", createCustomer(rubix, dbConn, logger))
	router.Get("/", getAllCustomers(dbConn, logger))
	router.Get("/unserved", getUnservedCustomers(dbConn, logger))
	router.Put("/", markAsServed(dbConn, logger))

	return router
}
