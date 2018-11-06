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

func createQueue(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var queue db.Queue
		err := json.NewDecoder(r.Body).Decode(&queue)
		if err != nil {
			handleServerError(w, "failed decoding request payload", err, logger)
			return
		}

		repo := db.NewQueuesRepo(dbConn)
		err = repo.Create(&queue)
		if err != nil {
			handleServerError(w, "failed creating queue", err, logger)
			return
		}

		render.JSON(w, r, Response{Info: "queue created successfully"})
	}
}

func getAllQueues(dbConn *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := db.NewQueuesRepo(dbConn)
		queues, err := repo.GetAll()
		if err != nil {
			handleServerError(w, "failed fetching all queues", err, logger)
			return
		}

		render.JSON(w, r, Response{Data: queues})
	}
}

func queuesRoutes(dbConn *sqlx.DB, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", createQueue(dbConn, logger))
	router.Get("/", getAllQueues(dbConn, logger))

	return router
}
