package main

import (
	"github.com/artrey/go-bank-service/pkg/transaction"
	"io"
	"log"
	"os"
)

func main() {
	svc := transaction.NewService()
	for i := 0; i < 100; i++ {
		svc.Add("4561 2612 1234 5467", "5106 2105 0000 0002", 100_00, 130_00, "0000")
	}

	if err := execute(svc, "transactions.csv"); err != nil {
		os.Exit(1)
	}
}

func execute(svc *transaction.Service, filename string) (err error) {
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

	err = svc.ExportAsCsv(file)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
