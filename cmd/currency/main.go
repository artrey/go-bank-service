package main

import (
	"context"
	"github.com/artrey/go-bank-service/pkg/currency"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}

func execute() error {
	const baseUrl = "https://raw.githubusercontent.com"
	const filename = "currencies.json"

	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Println(cerr)
			if err == nil {
				err = cerr
			}
		}
	}(file)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	svc := currency.NewService(baseUrl, http.DefaultClient)
	err = svc.Extract(ctx, file)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
