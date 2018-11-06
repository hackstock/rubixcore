package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// InitRoutes returns a http.Handler with all accessible endpoints registered
func InitRoutes(
	brokerConn *amqp.Connection,
	dbConn *sqlx.DB,
	upgrader *websocket.Upgrader,
	logger *zap.Logger,
) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		/*middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,*/
	)

	router.Mount("/users", usersRoutes(dbConn, logger))
	router.Mount("/queues", queuesRoutes(dbConn, logger))
	router.Mount("/customers", customersRoutes(dbConn, logger))

	return router
}
