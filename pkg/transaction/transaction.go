package transaction

import (
	"sort"
	"time"
)

type Transaction struct {
	Id        int64
	From      string
	To        string
	Timestamp int64
	Amount    int64
	Total     int64
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
