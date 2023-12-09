package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Gekuro/democrator/api/graph"
	"github.com/Gekuro/democrator/api/store"
	"github.com/joho/godotenv"
)

const DEFAULT_PORT = "8000";

func main() {
	// dotenv
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("error reading .env file: %s", err))
	}

	// store
	db, err := store.NewStore()
	if err != nil {
		panic(fmt.Errorf("error setting up store: %s", err))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(db)}))
	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
