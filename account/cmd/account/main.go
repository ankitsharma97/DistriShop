package main

import (
	"fmt"
	"log"
	"time"

	"microservice/account"
	envconfig "github.com/kelseyhightower/envconfig"
	retry"github.com/tinrab/retry"
)

type Config struct {
	databaseURL string `env:"DATABASE_URL,required"`
}


func main() {
	// Load configuration from environment variables.
	var cfg Config
	// err := env.Parse(&cfg)
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to parse env vars: %v", err)
	}

	// Initialize the PostgreSQL repository.
	var r account.Repository
	retry.ForeverSleep(5*time.Second, func(attempt int) error {
		var err error
		r, err = account.NewPostgresRepository(cfg.databaseURL)
		if err != nil {
			log.Printf("Failed to connect to Postgres (attempt %d): %v", attempt, err)
		}
		return err
	})
	defer r.Close()

	fmt.Println("Connected to Postgres")
	
	// Initialize the account service.
	svc := account.NewService(r)
	// Start the gRPC server.
	if err := account.ListenGRPC(svc, 8080); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}