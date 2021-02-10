package main

import (
	"bytes"
	"encoding/csv"
	"github.com/artrey/go-bank-service/pkg/transaction"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	svc := transaction.NewService()

	if err := execute(svc, "transactions.csv"); err != nil {
		os.Exit(1)
	}
}

func execute(svc *transaction.Service, filename string) (err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return err
	}

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		return err
	}

	svc.ImportRecords(records)
	return nil
}
