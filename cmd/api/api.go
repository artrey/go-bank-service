package main

import (
	"github.com/artrey/go-bank-service/pkg/app"
	"github.com/artrey/go-bank-service/pkg/card"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultHost = "localhost"
	defaultPort = "9999"
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

	address := net.JoinHostPort(host, port)
	log.Println(address)

	if err := execute(address); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(address string) error {
	cardSvc := card.NewService("MyBank", "5106 21")
	mux := http.NewServeMux()
	application := app.NewServer(cardSvc, mux)
	application.Init()

	server := http.Server{
		Addr:    address,
		Handler: application,
	}
	return server.ListenAndServe()
}
