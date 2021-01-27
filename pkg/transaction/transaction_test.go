package transaction_test

import (
	"github.com/artrey/go-bank-service/pkg/transaction"
	"reflect"
	"testing"
)

func TestService_Add(t *testing.T) {
	tests := []struct {
		name             string
		service          *transaction.Service
		wantTransactions int
	}{
		{
			name:             "empty service",
			service:          transaction.NewService(),
			wantTransactions: 1,
		},
		{
			name: "service with one transaction",
			service: &transaction.Service{
				&transaction.Transaction{
					Id:     1,
					From:   "4561 2612 1234 5464",
					To:     "2612 4561 1234 5464",
					Amount: 100_00,
					Total:  100_00,
				},
			},
			wantTransactions: 2,
		},
	}

	for _, tt := range tests {
		tt.service.Add(
			"4561 2612 1234 5464",
			"2612 4561 1234 5464",
			100_00,
			100_00,
		)
		gotTransactions := tt.service.Count()
		if gotTransactions != tt.wantTransactions {
			t.Errorf("%v: got = %v, want %v", tt.name, gotTransactions, tt.wantTransactions)
		}
	}
}

func TestSort(t *testing.T) {
	transactions := []*transaction.Transaction{
		{
			Amount: 100_00,
			Total:  100_00,
		},
		{
			Amount: 200_00,
			Total:  200_00,
		},
		{
			Amount: 300_00,
			Total:  300_00,
		},
	}

	tests := []struct {
		name         string
		transactions []*transaction.Transaction
		wantOrder    []*transaction.Transaction
	}{
		{
			name: "ordered",
			transactions: []*transaction.Transaction{
				transactions[2],
				transactions[1],
				transactions[0],
			},
			wantOrder: []*transaction.Transaction{
				transactions[2],
				transactions[1],
				transactions[0],
			},
		},
		{
			name: "unordered",
			transactions: []*transaction.Transaction{
				transactions[1],
				transactions[2],
				transactions[0],
			},
			wantOrder: []*transaction.Transaction{
				transactions[2],
				transactions[1],
				transactions[0],
			},
		},
	}

	for _, tt := range tests {
		if got := transaction.Sort(tt.transactions); !reflect.DeepEqual(got, tt.wantOrder) {
			t.Errorf("%v: got = %v, want %v", tt.name, got, tt.wantOrder)
		}
	}
}
