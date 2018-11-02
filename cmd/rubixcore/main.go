package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
	"github.com/hackstock/rubixcore/pkg/api"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const (
	production  = "production"
	development = "development"
)

var env = struct {
	Port             int    `envconfig:"PORT" required:"true"`
	Environment      string `envconfig:"ENVIRONMENT" default:"development"`
	TicketsResetTime string `envconfig:"TICKETS_RESET_TIME" required:"true"`
	ServiceDSN       string `envconfig:"SERVICE_DSN" required:"true"`
	RabbitMQURL      string `envconfig:"RABBITMQ_URL" required:"true"`
	JWTIssuer        string `envconfig:"JWT_ISSUER" required:"true"`
	JWTSecret        string `envconfig:"JWT_SECRET" required:"true"`
}{}

func init() {
	err := envconfig.Process("", &env)
	if err != nil {
		failOnError("failed loading configurations", err)
	}
}

func initLogger(environment string) (*zap.Logger, error) {
	if environment == production {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}

func main() {
	logger, err := initLogger(env.Environment)
	failOnError("failed initializing logger", err)

	if env.Environment == development {
		logger.Info("configurations loaded successfully", zap.Any("configs", env))
	}

	brokerConn, err := amqp.Dial(env.RabbitMQURL)
	failOnError("failed connecting to rabbitmq", err)

	logger.Info("connected to rabbitmq successfully")

	dbConn, err := sqlx.Open("mysql", env.ServiceDSN)
	failOnError("failed connecting to mysql", err)

	logger.Info("connected to mysql successfully")

	dbConn.SetConnMaxLifetime(time.Second * 14400)
	dbConn.SetMaxIdleConns(50)
	dbConn.SetMaxOpenConns(100)

	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		logger.Fatal("failed binding to port", zap.Int("port", env.Port))
	}
	defer listener.Close()

	url := fmt.Sprintf("http://%s", listener.Addr())
	logger.Info("server listening on ", zap.String("url", url))

	router := api.InitRoutes(
		brokerConn, 
		dbConn,
		&websocket.Upgrader{},
		 logger,
	)

	server := &http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router),
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	idleConnsClosed := make(chan struct{})
	go func() {
		defer close(idleConnsClosed)

		recv := <-sigs
		logger.Info("received signal, shutting down", zap.Any("signal", recv.String))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Warn("error shutting down server", zap.Error(err))
		}
	}()

	if err = server.Serve(listener); err != nil {
		if err != http.ErrServerClosed {
			logger.Fatal("server returned error", zap.Error(err))
		}
	}

	<-idleConnsClosed
	logger.Info("server shutdown successfully")
}

func failOnError(msg string, err error) {
	if err != nil {
		log.Fatalf("%s : %v", msg, err)
	}
}
