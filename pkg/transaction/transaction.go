package transaction

import (
	"encoding/csv"
	"errors"
	"github.com/artrey/go-bank-service/pkg/mcc"
	"io"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"
)

var (
	InvalidSizeCsvSlice = errors.New("size of slice with data")
)

type Transaction struct {
	XMLName   string  `json:"-" xml:"transaction"`
	Id        int64   `json:"id" xml:"id"`
	From      string  `json:"from" xml:"from"`
	To        string  `json:"to" xml:"to"`
	Timestamp int64   `json:"timestamp" xml:"timestamp"`
	Amount    int64   `json:"amount" xml:"amount"`
	Total     int64   `json:"total" xml:"total"`
	MCC       mcc.MCC `json:"mcc" xml:"mcc"`
}

type Transactions struct {
	XMLName      string         `xml:"transactions"`
	Transactions []*Transaction `xml:"transaction"`
}

type Service struct {
	mu           sync.Mutex
	transactions []*Transaction
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Count() int {
	return len(s.transactions)
}

func (s *Service) Add(from, to string, amount, total int64, MCC mcc.MCC) *Transaction {
	id := int64(1)

	s.mu.Lock()

	if s.Count() > 0 {
		id = s.transactions[s.Count()-1].Id + 1
	}
	transaction := &Transaction{
		Id:        id,
		From:      from,
		To:        to,
		Timestamp: time.Now().UTC().Unix(),
		Amount:    amount,
		Total:     total,
		MCC:       MCC,
	}
	s.transactions = append(s.transactions, transaction)

	s.mu.Unlock()

	return transaction
}

func (s *Service) Transactions() []*Transaction {
	s.mu.Lock()
	data := s.transactions[:]
	s.mu.Unlock()
	return data
}

func (s *Service) ExportRecords() [][]string {
	s.mu.Lock()
	if s.Count() == 0 {
		s.mu.Unlock()
		return [][]string{}
	}

	records := make([][]string, 0, s.Count())
	for _, t := range s.transactions {
		record := AsCsvSlice(t)
		records = append(records, record)
	}
	s.mu.Unlock()

	return records
}

func (s *Service) ExportAsCsv(writer io.Writer) error {
	records := s.ExportRecords()
	if len(records) == 0 {
		return nil
	}
	w := csv.NewWriter(writer)
	return w.WriteAll(records)
}

func (s *Service) ImportRecords(records [][]string) {
	if len(records) == 0 {
		return
	}

	s.mu.Lock()
	for line, record := range records {
		t, err := FromCsvSlice(record)
		if err != nil {
			log.Printf("line %d | %v", line+1, err)
			continue
		}
		s.transactions = append(s.transactions, t)
	}
	s.mu.Unlock()
}

func AsCsvSlice(t *Transaction) []string {
	return []string{
		strconv.FormatInt(t.Id, 10),
		t.From,
		t.To,
		strconv.FormatFloat(float64(t.Amount)/100.0, 'f', 2, 64),
		strconv.FormatFloat(float64(t.Total)/100.0, 'f', 2, 64),
		strconv.FormatInt(t.Timestamp, 10),
		string(t.MCC),
	}
}

func FromCsvSlice(data []string) (t *Transaction, err error) {
	if len(data) != 7 {
		return nil, InvalidSizeCsvSlice
	}

	id, err := strconv.ParseInt(data[0], 10, 64)
	if err != nil {
		return nil, err
	}

	amountFloat, err := strconv.ParseFloat(data[3], 64)
	if err != nil {
		return nil, err
	}
	amount := int64(amountFloat * 100)

	totalFloat, err := strconv.ParseFloat(data[4], 64)
	if err != nil {
		return nil, err
	}
	total := int64(totalFloat * 100)

	timestamp, err := strconv.ParseInt(data[5], 10, 64)
	if err != nil {
		return nil, err
	}

	t = &Transaction{
		Id:        id,
		From:      data[1],
		To:        data[2],
		Timestamp: timestamp,
		Amount:    amount,
		Total:     total,
		MCC:       mcc.MCC(data[6]),
	}
	return t, nil
}

func Sort(transactions []*Transaction) []*Transaction {
	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].Total > transactions[j].Total
	})
	return transactions
}

func Slice(transactions []*Transaction, from, to int64) []*Transaction {
	validIndices := make([]int, 0)
	for idx, transaction := range transactions {
		if from <= transaction.Timestamp && transaction.Timestamp < to {
			validIndices = append(validIndices, idx)
		}
	}

	if len(validIndices) == 0 {
		return []*Transaction{}
	}

	return transactions[validIndices[0] : validIndices[len(validIndices)-1]+1]
}

func Sum(transactions []*Transaction) int64 {
	result := int64(0)
	for _, transaction := range transactions {
		result += transaction.Total
	}
	return result
}

func SumByMCC(transactions []*Transaction) map[mcc.MCC]int64 {
	result := make(map[mcc.MCC]int64)
	for _, transaction := range transactions {
		result[transaction.MCC] += transaction.Total
	}
	return result
}

func FilterByCardNumber(transactions []*Transaction, fromCardNumber string) []*Transaction {
	cardTransactions := make([]*Transaction, 0)
	for _, transaction := range transactions {
		if transaction.From == fromCardNumber {
			cardTransactions = append(cardTransactions, transaction)
		}
	}
	return cardTransactions
}

func Categorize(transactions []*Transaction, fromCardNumber string) map[string]int64 {
	cardTransactions := FilterByCardNumber(transactions, fromCardNumber)
	sums := SumByMCC(cardTransactions)
	result := make(map[string]int64)
	for MCC, sum := range sums {
		result[MCC.ToCategory()] += sum
	}
	return result
}

func CategorizeConcurrentWithMutex(transactions []*Transaction,
	fromCardNumber string, goroutines int) map[string]int64 {

	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	mu := sync.Mutex{}
	result := make(map[string]int64)

	partSize := len(transactions) / goroutines
	for n := 0; n < goroutines; n++ {
		fromIndex := n * partSize
		toIndex := (n + 1) * partSize
		if n == goroutines-1 {
			toIndex = len(transactions)
		}
		part := transactions[fromIndex:toIndex]

		go func() {
			m := Categorize(part, fromCardNumber)

			mu.Lock()
			for category, sum := range m {
				result[category] += sum
			}
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()

	return result
}

func CategorizeConcurrentWithChannels(transactions []*Transaction,
	fromCardNumber string, goroutines int) map[string]int64 {

	result := make(map[string]int64)
	ch := make(chan map[string]int64)

	partSize := len(transactions) / goroutines
	for n := 0; n < goroutines; n++ {
		fromIndex := n * partSize
		toIndex := (n + 1) * partSize
		if n == goroutines-1 {
			toIndex = len(transactions)
		}
		part := transactions[fromIndex:toIndex]

		go func(ch chan<- map[string]int64) {
			ch <- Categorize(part, fromCardNumber)
		}(ch)
	}

	finished := 0
	for value := range ch {
		for category, sum := range value {
			result[category] += sum
		}

		finished++
		if finished == goroutines {
			close(ch)
			break
		}
	}

	return result
}

func CategorizeConcurrentWithMutexManual(transactions []*Transaction,
	fromCardNumber string, goroutines int) map[string]int64 {

	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	mu := sync.Mutex{}
	result := make(map[string]int64)

	partSize := len(transactions) / goroutines
	for n := 0; n < goroutines; n++ {
		fromIndex := n * partSize
		toIndex := (n + 1) * partSize
		if n == goroutines-1 {
			toIndex = len(transactions)
		}
		part := transactions[fromIndex:toIndex]

		go func() {
			for _, transaction := range part {
				if transaction.From == fromCardNumber {
					mu.Lock()
					result[transaction.MCC.ToCategory()] += transaction.Total
					mu.Unlock()
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()

	return result
}
