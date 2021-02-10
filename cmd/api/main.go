package main

import (
	"github.com/artrey/go-bank-service/pkg/app"
	"github.com/artrey/go-bank-service/pkg/card"
	"net/http"
	"os"
)

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}

func execute() error {
	cardSvc := card.NewService("MyBank", "5106 21")
	mux := http.NewServeMux()
	application := app.NewServer(cardSvc, mux)
	application.Init()

	server := http.Server{
		Addr:    "0.0.0.0:9999",
		Handler: application,
	}
	return server.ListenAndServe()
}
