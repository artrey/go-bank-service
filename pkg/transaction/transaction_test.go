package transaction_test

import (
	"github.com/artrey/go-bank-service/pkg/transaction"
	"reflect"
	"testing"
)

func TestService_Add(t *testing.T) {
	preparedService := transaction.NewService()
	preparedService.Add("4561 2612 1234 5464", "2612 4561 1234 5464",
		100_00, 100_00, "0000")

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
			name:             "service with one transaction",
			service:          preparedService,
			wantTransactions: 2,
		},
	}

	for _, tt := range tests {
		tt.service.Add(
			"4561 2612 1234 5464",
			"2612 4561 1234 5464",
			100_00,
			100_00,
			"2222",
		)
		gotTransactions := tt.service.Count()
		if gotTransactions != tt.wantTransactions {
			t.Errorf("%v: got = %v, want %v", tt.name, gotTransactions, tt.wantTransactions)
		}
	}
}

func TestFromCsvSlice(t *testing.T) {
	tests := []struct {
		name       string
		data       []string
		wantResult *transaction.Transaction
		wantError  error
	}{
		{
			name: "valid",
			data: []string{
				"14",
				"4561 2612 1234 5467",
				"5106 2105 0000 0002",
				"100.00",
				"130.00",
				"1611837340",
				"0000",
			},
			wantResult: &transaction.Transaction{
				Id:        14,
				From:      "4561 2612 1234 5467",
				To:        "5106 2105 0000 0002",
				Timestamp: 1611837340,
				Amount:    100_00,
				Total:     130_00,
				MCC:       "0000",
			},
			wantError: nil,
		},
		{
			name: "invalid size of slice",
			data: []string{
				"14",
				"4561 2612 1234 5467",
				"100.00",
				"130.00",
			},
			wantResult: nil,
			wantError:  transaction.InvalidSizeCsvSlice,
		},
	}

	for _, tt := range tests {
		gotResult, gotError := transaction.FromCsvSlice(tt.data)
		if gotError != tt.wantError {
			t.Errorf("%v (error): got = %v, want %v", tt.name, gotError, tt.wantError)
		}
		if !reflect.DeepEqual(gotResult, tt.wantResult) {
			t.Errorf("%v (result): got = %v, want %v", tt.name, gotResult, tt.wantResult)
		}
	}
}

func TestSort(t *testing.T) {
	transactions := []*transaction.Transaction{
		{
			Total: 100_00,
		},
		{
			Total: 200_00,
		},
		{
			Total: 300_00,
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

func makeExpectedResult() map[string]int64 {
	return map[string]int64{
		"Супермаркеты":         10_000_00,
		"Автоуслуги":           20_000_00,
		"Рестораны":            10_000_00,
		"Аптеки":               10_000_00,
		"Категория не указана": 10_000_00,
	}
}

type categorizeArgs struct {
	transactions   []*transaction.Transaction
	fromCardNumber string
	goroutines     int
}

func TestCategorize(t *testing.T) {
	tests := []struct {
		name     string
		args     categorizeArgs
		expected map[string]int64
	}{
		{
			name: "Positive categorizing",
			args: categorizeArgs{
				transactions:   makeTransactions(),
				fromCardNumber: "4561 2612 1234 5467",
			},
			expected: makeExpectedResult(),
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
	tests := []struct {
		name     string
		args     categorizeArgs
		expected map[string]int64
	}{
		{
			name: "Positive categorizing",
			args: categorizeArgs{
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
	tests := []struct {
		name     string
		args     categorizeArgs
		expected map[string]int64
	}{
		{
			name: "Positive categorizing",
			args: categorizeArgs{
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
	tests := []struct {
		name     string
		args     categorizeArgs
		expected map[string]int64
	}{
		{
			name: "Positive categorizing",
			args: categorizeArgs{
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

func BenchmarkCategorize(b *testing.B) {
	args := categorizeArgs{
		transactions:   makeTransactions(),
		fromCardNumber: "4561 2612 1234 5467",
	}
	expected := makeExpectedResult()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := transaction.Categorize(args.transactions, args.fromCardNumber)
		b.StopTimer()
		if !reflect.DeepEqual(got, expected) {
			b.Fatalf("invalid result, got = %v, want %v", got, expected)
		}
		b.StartTimer()
	}
}

func BenchmarkCategorizeConcurrentWithMutex(b *testing.B) {
	args := categorizeArgs{
		transactions:   makeTransactions(),
		fromCardNumber: "4561 2612 1234 5467",
		goroutines:     10,
	}
	expected := makeExpectedResult()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := transaction.CategorizeConcurrentWithMutex(args.transactions, args.fromCardNumber, args.goroutines)
		b.StopTimer()
		if !reflect.DeepEqual(got, expected) {
			b.Fatalf("invalid result, got = %v, want %v", got, expected)
		}
		b.StartTimer()
	}
}

func BenchmarkCategorizeConcurrentWithMutex100g(b *testing.B) {
	args := categorizeArgs{
		transactions:   makeTransactions(),
		fromCardNumber: "4561 2612 1234 5467",
		goroutines:     100,
	}
	expected := makeExpectedResult()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := transaction.CategorizeConcurrentWithMutex(args.transactions, args.fromCardNumber, args.goroutines)
		b.StopTimer()
		if !reflect.DeepEqual(got, expected) {
			b.Fatalf("invalid result, got = %v, want %v", got, expected)
		}
		b.StartTimer()
	}
}

func BenchmarkCategorizeConcurrentWithMutex1000g(b *testing.B) {
	args := categorizeArgs{
		transactions:   makeTransactions(),
		fromCardNumber: "4561 2612 1234 5467",
		goroutines:     1000,
	}
	expected := makeExpectedResult()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := transaction.CategorizeConcurrentWithMutex(args.transactions, args.fromCardNumber, args.goroutines)
		b.StopTimer()
		if !reflect.DeepEqual(got, expected) {
			b.Fatalf("invalid result, got = %v, want %v", got, expected)
		}
		b.StartTimer()
	}
}

func BenchmarkCategorizeConcurrentWithChannels(b *testing.B) {
	args := categorizeArgs{
		transactions:   makeTransactions(),
		fromCardNumber: "4561 2612 1234 5467",
		goroutines:     10,
	}
	expected := makeExpectedResult()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := transaction.CategorizeConcurrentWithChannels(args.transactions, args.fromCardNumber, args.goroutines)
		b.StopTimer()
		if !reflect.DeepEqual(got, expected) {
			b.Fatalf("invalid result, got = %v, want %v", got, expected)
		}
		b.StartTimer()
	}
}

func BenchmarkCategorizeConcurrentWithChannels100g(b *testing.B) {
	args := categorizeArgs{
		transactions:   makeTransactions(),
		fromCardNumber: "4561 2612 1234 5467",
		goroutines:     100,
	}
	expected := makeExpectedResult()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := transaction.CategorizeConcurrentWithChannels(args.transactions, args.fromCardNumber, args.goroutines)
		b.StopTimer()
		if !reflect.DeepEqual(got, expected) {
			b.Fatalf("invalid result, got = %v, want %v", got, expected)
		}
		b.StartTimer()
	}
}

func BenchmarkCategorizeConcurrentWithChannels1000g(b *testing.B) {
	args := categorizeArgs{
		transactions:   makeTransactions(),
		fromCardNumber: "4561 2612 1234 5467",
		goroutines:     1000,
	}
	expected := makeExpectedResult()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := transaction.CategorizeConcurrentWithChannels(args.transactions, args.fromCardNumber, args.goroutines)
		b.StopTimer()
		if !reflect.DeepEqual(got, expected) {
			b.Fatalf("invalid result, got = %v, want %v", got, expected)
		}
		b.StartTimer()
	}
}

func BenchmarkCategorizeConcurrentWithMutexManual(b *testing.B) {
	args := categorizeArgs{
		transactions:   makeTransactions(),
		fromCardNumber: "4561 2612 1234 5467",
		goroutines:     10,
	}
	expected := makeExpectedResult()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := transaction.CategorizeConcurrentWithMutexManual(args.transactions, args.fromCardNumber, args.goroutines)
		b.StopTimer()
		if !reflect.DeepEqual(got, expected) {
			b.Fatalf("invalid result, got = %v, want %v", got, expected)
		}
		b.StartTimer()
	}
}

func BenchmarkCategorizeConcurrentWithMutexManual100g(b *testing.B) {
	args := categorizeArgs{
		transactions:   makeTransactions(),
		fromCardNumber: "4561 2612 1234 5467",
		goroutines:     100,
	}
	expected := makeExpectedResult()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := transaction.CategorizeConcurrentWithMutexManual(args.transactions, args.fromCardNumber, args.goroutines)
		b.StopTimer()
		if !reflect.DeepEqual(got, expected) {
			b.Fatalf("invalid result, got = %v, want %v", got, expected)
		}
		b.StartTimer()
	}
}

func BenchmarkCategorizeConcurrentWithMutexManual1000g(b *testing.B) {
	args := categorizeArgs{
		transactions:   makeTransactions(),
		fromCardNumber: "4561 2612 1234 5467",
		goroutines:     1000,
	}
	expected := makeExpectedResult()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := transaction.CategorizeConcurrentWithMutexManual(args.transactions, args.fromCardNumber, args.goroutines)
		b.StopTimer()
		if !reflect.DeepEqual(got, expected) {
			b.Fatalf("invalid result, got = %v, want %v", got, expected)
		}
		b.StartTimer()
	}
}
