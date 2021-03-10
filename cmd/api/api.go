package main

import (
	"context"
	"github.com/artrey/go-bank-service/pkg/storage/postgres"
	"log"
	"net"
	"os"
)

const (
	defaultHost = "0.0.0.0"
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
	//cardSvc := card.NewService("MyBank", "5106 21")
	//mux := http.NewServeMux()
	//application := app.NewServer(cardSvc, mux)
	//application.Init()
	//
	//server := http.Server{
	//	Addr:    address,
	//	Handler: application,
	//}
	//return server.ListenAndServe()

	storage, err := postgres.New(context.Background(), "postgres://go:go@172.30.235.58:5532/go")
	if err != nil {
		log.Println(err)
		return err
	}

	cards, err := storage.GetCardsByClientId(2)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, c := range cards {
		log.Printf("%+v", c)
		transactions, err := storage.GetTransactionsByCardId(c.Id)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, t := range transactions {
			log.Printf("%+v", t)
		}
	}

	return nil
}
