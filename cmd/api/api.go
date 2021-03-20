package main

import (
	"context"
	"github.com/artrey/go-bank-service/pkg/service"
	"github.com/artrey/go-bank-service/pkg/storage/postgres"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = "9999"
	defaultDsn  = "postgres://user:pass@localhost:5432/api"
)

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}
	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		dsn = defaultDsn
	}
	log.Println(dsn)

	address := net.JoinHostPort(host, port)
	log.Println(address)

	if err := execute(address, dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(address, dsn string) error {
	storage, err := postgres.New(context.Background(), dsn)
	if err != nil {
		log.Println(err)
		return err
	}
	defer storage.Close()

	mux := http.NewServeMux()
	application := service.NewServer(storage, mux)
	application.Init()

	server := http.Server{
		Addr:    address,
		Handler: application,
	}
	return server.ListenAndServe()
}
