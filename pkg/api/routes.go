package api

import (
	"github.com/go-chi/chi"
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

	return router
}
