package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Gekuro/democrator/api/graph"
	"github.com/Gekuro/democrator/api/logger"
	"github.com/Gekuro/democrator/api/store"
	"github.com/joho/godotenv"
)

func main() {
	// dotenv
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error reading .env file: %s", err)
	}

	// store
	db, err := store.NewStore()
	if err != nil {
		log.Fatalf("error setting up store: %s", err)
	}

	// logger
	logWriter, err := logger.NewLoggerWriter()
	if err != nil {
		log.Fatalf("error setting up store: %s", err)
	}
	log.SetOutput(logWriter)

	// gql
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(db)}))
	srv.AddTransport(&transport.Websocket{})

	if strings.ToUpper(os.Getenv("APP_ENV")) == "DEV" {
		http.Handle("/", playground.Handler("Democrator GraphiQL", "/query"))
	}
	http.Handle("/query", srv)

	fmt.Printf("connect to http://localhost:%s/ for GraphiQL", os.Getenv("PORT")) // intentionally omitted from logs
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
