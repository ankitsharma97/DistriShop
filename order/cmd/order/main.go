package main

import (
	"log"
	"microservice/order"
	"time"

	githubenv "github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
}

func main() {
	var cfg Config
	// Load environment variables into cfg
	if err := githubenv.Process("", &cfg); err != nil {
		panic(err)
	}
	// if cfg.DatabaseURL == "" {
	// 	panic("DATABASE_URL is required")
	// }

	var r order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println("retry connecting ", err)
		}
		return
	})
	defer r.Close()
	log.Println("Listening on port 8080...")
	s := order.NewOrderService(r)
	log.Fatal(order.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, 8080))
}
