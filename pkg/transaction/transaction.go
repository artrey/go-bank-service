package transaction

import (
	"github.com/artrey/go-bank-service/pkg/mcc"
	"sort"
	"sync"
	"time"
)

type Transaction struct {
	Id        int64
	From      string
	To        string
	Timestamp int64
	Amount    int64
	Total     int64
	MCC       mcc.MCC
}

type Service []*Transaction

func NewService() *Service {
	return &Service{}
}

func (s *Service) Count() int {
	return len(*s)
}

func (s *Service) Add(from, to string, amount, total int64) *Transaction {
	var id int64 = 1
	if len(*s) > 0 {
		id = (*s)[len(*s)-1].Id + 1
	}
	transaction := &Transaction{
		Id:        id,
		From:      from,
		To:        to,
		Timestamp: time.Now().UTC().Unix(),
		Amount:    amount,
		Total:     total,
	}
	*s = append(*s, transaction)
	return transaction
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

		go func(ch chan <- map[string]int64) {
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
