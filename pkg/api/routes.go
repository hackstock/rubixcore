package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// InitRoutes returns a http.Handler with all accessible endpoints registered
func InitRoutes(upgrader *websocket.Upgrader, logger *zap.Logger) http.Handler {
	router := chi.NewRouter()

	return router
}
