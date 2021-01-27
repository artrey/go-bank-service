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

func makeTransactions() []*transaction.Transaction {
	const users = 1_000
	const transactionsPerUser = 1_000
	const transactionTotal = 1_00
	transactions := make([]*transaction.Transaction, users*transactionsPerUser)
	for index := range transactions {
		switch index % 100 {
		case 0:
			transactions[index] = &transaction.Transaction{
				From:  "4561 2612 1234 5467",
				Total: transactionTotal,
				MCC:   "5411",
			}
		case 10:
			transactions[index] = &transaction.Transaction{
				From:  "4561 2612 1234 5467",
				Total: transactionTotal,
				MCC:   "5812",
			}
		case 20:
			transactions[index] = &transaction.Transaction{
				From:  "4561 2612 1234 5467",
				Total: transactionTotal,
				MCC:   "5912",
			}
		case 30, 31:
			transactions[index] = &transaction.Transaction{
				From:  "4561 2612 1234 5467",
				Total: transactionTotal,
				MCC:   "5533",
			}
		case 40:
			transactions[index] = &transaction.Transaction{
				From:  "4561 2612 1234 5467",
				Total: transactionTotal,
				MCC:   "2222",
			}
		default:
			transactions[index] = &transaction.Transaction{
				From:  "1234 4561 2612 5467",
				Total: transactionTotal,
				MCC:   "5533",
			}
		}
	}
	return transactions
}

func TestCategorize(t *testing.T) {
	type args struct {
		transactions   []*transaction.Transaction
		fromCardNumber string
	}
	tests := []struct {
		name     string
		args     args
		expected map[string]int64
	}{
		{
			name: "Positive categorizing",
			args: args{
				transactions:   makeTransactions(),
				fromCardNumber: "4561 2612 1234 5467",
			},
			expected: map[string]int64{
				"Супермаркеты":         10_000_00,
				"Автоуслуги":           20_000_00,
				"Рестораны":            10_000_00,
				"Аптеки":               10_000_00,
				"Категория не указана": 10_000_00,
			},
		},
	}

	for _, tt := range tests {
		got := transaction.Categorize(tt.args.transactions, tt.args.fromCardNumber)
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("%v: got = %v, want %v", tt.name, got, tt.expected)
		}
	}
}

func TestCategorizeConcurrentWithMutex(t *testing.T) {
	type args struct {
		transactions   []*transaction.Transaction
		fromCardNumber string
		goroutines     int
	}
	tests := []struct {
		name     string
		args     args
		expected map[string]int64
	}{
		{
			name: "Positive categorizing",
			args: args{
				transactions:   makeTransactions(),
				fromCardNumber: "4561 2612 1234 5467",
				goroutines:     100,
			},
			expected: map[string]int64{
				"Супермаркеты":         10_000_00,
				"Автоуслуги":           20_000_00,
				"Рестораны":            10_000_00,
				"Аптеки":               10_000_00,
				"Категория не указана": 10_000_00,
			},
		},
	}

	for _, tt := range tests {
		got := transaction.CategorizeConcurrentWithMutex(
			tt.args.transactions, tt.args.fromCardNumber, tt.args.goroutines)
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("%v: got = %v, want %v", tt.name, got, tt.expected)
		}
	}
}

func TestCategorizeConcurrentWithChannels(t *testing.T) {
	type args struct {
		transactions   []*transaction.Transaction
		fromCardNumber string
		goroutines     int
	}
	tests := []struct {
		name     string
		args     args
		expected map[string]int64
	}{
		{
			name: "Positive categorizing",
			args: args{
				transactions:   makeTransactions(),
				fromCardNumber: "4561 2612 1234 5467",
				goroutines:     10,
			},
			expected: map[string]int64{
				"Супермаркеты":         10_000_00,
				"Автоуслуги":           20_000_00,
				"Рестораны":            10_000_00,
				"Аптеки":               10_000_00,
				"Категория не указана": 10_000_00,
			},
		},
	}

	for _, tt := range tests {
		got := transaction.CategorizeConcurrentWithChannels(
			tt.args.transactions, tt.args.fromCardNumber, tt.args.goroutines)
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("%v: got = %v, want %v", tt.name, got, tt.expected)
		}
	}
}

func TestCategorizeConcurrentWithMutexManual(t *testing.T) {
	type args struct {
		transactions   []*transaction.Transaction
		fromCardNumber string
		goroutines     int
	}
	tests := []struct {
		name     string
		args     args
		expected map[string]int64
	}{
		{
			name: "Positive categorizing",
			args: args{
				transactions:   makeTransactions(),
				fromCardNumber: "4561 2612 1234 5467",
				goroutines:     100,
			},
			expected: map[string]int64{
				"Супермаркеты":         10_000_00,
				"Автоуслуги":           20_000_00,
				"Рестораны":            10_000_00,
				"Аптеки":               10_000_00,
				"Категория не указана": 10_000_00,
			},
		},
	}

	for _, tt := range tests {
		got := transaction.CategorizeConcurrentWithMutexManual(
			tt.args.transactions, tt.args.fromCardNumber, tt.args.goroutines)
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("%v: got = %v, want %v", tt.name, got, tt.expected)
		}
	}
}
