package main

import (
	"encoding/json"
	"fmt"
	"github.com/artrey/go-bank-service/pkg/transaction"
	"log"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}()
	err = trace.Start(f)
	if err != nil {
		log.Fatal(err)
	}
	defer trace.Stop()

	transactions := []*transaction.Transaction{
		{
			From:      "4561 2612 1234 5467",
			Timestamp: time.Date(2020, 9, 25, 16, 0, 0, 0, time.Local).UTC().Unix(),
			Total:     500_00,
			MCC:       "5411",
		},
		{
			From:      "4561 2612 1234 5467",
			Timestamp: time.Date(2020, 9, 26, 12, 0, 0, 0, time.Local).UTC().Unix(),
			Total:     500_00,
		},
		{
			From:      "4561 2612 1234 5467",
			Timestamp: time.Date(2020, 10, 4, 20, 15, 0, 0, time.Local).UTC().Unix(),
			Total:     1200_00,
		},
		{
			From:      "4561 2612 1234 5467",
			Timestamp: time.Date(2021, 1, 22, 20, 15, 0, 0, time.Local).UTC().Unix(),
			Total:     100_00,
		},
		{
			From:      "4561 2612 1234 5467",
			Timestamp: time.Date(2021, 1, 23, 23, 59, 59, 0, time.Local).UTC().Unix(),
			Total:     15000_00,
		},
	}
	from := time.Date(2020, 9, 15, 0, 0, 0, 0, time.Local)
	to := time.Date(2021, 1, 25, 0, 0, 0, 0, time.Local)
	SumConcurrently(transactions, from.UTC().Unix(), to.UTC().Unix())

	transaction.CategorizeConcurrentWithMutex(transactions, "4561 2612 1234 5467", 5)
	transaction.CategorizeConcurrentWithChannels(transactions, "4561 2612 1234 5467", 5)

	encoded, err := json.Marshal(transactions)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(string(encoded))

	var decoded []*transaction.Transaction
	err = json.Unmarshal(encoded, &decoded)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	for _, t := range decoded {
		log.Printf("%#v", t)
	}
}

func SumConcurrently(transactions []*transaction.Transaction, from, to int64) int64 {
	fromDate := time.Unix(from, 0).Local()
	currentDate := time.Date(fromDate.Year(), fromDate.Month()+1, 1, 0, 0, 0, 0, time.Local)

	partsTimestamps := []int64{from}
	for currentDate.UTC().Unix() < to {
		partsTimestamps = append(partsTimestamps, currentDate.UTC().Unix())
		currentDate = currentDate.AddDate(0, 1, 0)
	}
	partsTimestamps = append(partsTimestamps, to)

	goroutines := len(partsTimestamps) - 1
	total := int64(0)

	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		fromPart := partsTimestamps[i]
		toPart := partsTimestamps[i+1]
		part := transaction.Slice(transactions, fromPart, toPart)
		go func() {
			sum := transaction.Sum(part)
			fmt.Println(
				"from", time.Unix(fromPart, 0).Local(),
				"to", time.Unix(toPart, 0).Local(),
				":", sum)
			wg.Done()
		}()
	}

	wg.Wait()
	return total
}
