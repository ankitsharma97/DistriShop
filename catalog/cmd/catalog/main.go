package main

import (
	"log"
	"microservice/catalog"
	"time"

	githubenv "github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
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

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println("retry connecting elastic:", err)
		}
		return
	})
	defer r.Close()
	log.Println("Listening on port 8080...")
	s := catalog.NewCatalogService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
