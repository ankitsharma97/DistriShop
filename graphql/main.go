package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

// AppConfig holds service endpoint configuration. Default values allow running the
// gateway standalone with in-memory resolvers (see graph.go). Replace the URLs
// when real microservices become available.
type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_URL" default:"http://localhost:4001"`
	CatalogURL string `envconfig:"CATALOG_URL" default:"http://localhost:4002"`
	OrderURL   string `envconfig:"ORDER_URL" default:"http://localhost:4003"`
}

func main() {
	var config AppConfig
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	// Build GraphQL server (currently in-memory data store)
	srv, err := NewGraphQLServer(config.AccountURL, config.CatalogURL, config.OrderURL)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()

	// Create the gqlgen HTTP handler
	gqlHandler := handler.NewDefaultServer(srv.ToExecutableSchema())

	http.Handle("/graphql", gqlHandler)
	// Enable GraphQL Playground at /playground for interactive queries
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	log.Println("GraphQL server running on :8080 (playground at /playground)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
