package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
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
}

func failOnError(msg string, err error) {
	if err != nil {
		log.Fatalf("%s : %v", msg, err)
	}
}
