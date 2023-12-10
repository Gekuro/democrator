package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Gekuro/democrator/api/graph"
	"github.com/Gekuro/democrator/api/store"
	"github.com/joho/godotenv"
)

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

	// logger
	logFile, err := os.OpenFile("server.log", os.O_APPEND, 0777)
	if err != nil {
		panic(fmt.Errorf("error openning log file: %s", err))
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// gql
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(db)}))
	srv.AddTransport(&transport.Websocket{})

	if strings.ToUpper(os.Getenv("APP_ENV")) == "DEV" {
		http.Handle("/", playground.Handler("Democrator GraphiQL", "/query"))
	}
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphiQL", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
