package transfer_test

import (
	"github.com/artrey/go-bank-service/pkg/card"
	"github.com/artrey/go-bank-service/pkg/transaction"
	"github.com/artrey/go-bank-service/pkg/transfer"
	"math"
	"testing"
)

func TestService_IsValidCardNumber(t *testing.T) {
	tests := []struct {
		name   string
		number string
		wantOk bool
	}{
		{
			name:   "Empty card number",
			number: "",
			wantOk: false,
		},
		{
			name:   "Valid card number",
			number: "4561 2612 1234 5467",
			wantOk: true,
		},
		{
			name:   "Invalid card number",
			number: "4561 2612 1234 5464",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		gotOk := transfer.IsValidCardNumber(tt.number)
		if gotOk != tt.wantOk {
			t.Errorf("%v: gotTotal = %v, want %v", tt.name, gotOk, tt.wantOk)
		}
	}
}

func makeEmptyCardService() *card.Service {
	return card.NewService("Tinkoff", "5106 21")
}

func makeCardService1Card() *card.Service {
	svc := makeEmptyCardService()
	svc.Issue(nil, "Visa", 10, card.Plastic, 1000_00,
		"RUB", "5106 2107 0000 0000", "https://...")
	return svc
}

func makeCardService2Cards() *card.Service {
	svc := makeCardService1Card()
	svc.Issue(nil, "MasterCard", 10, card.Plastic, 1000_00,
		"RUB", "5106 2105 0000 0002", "https://...")
	return svc
}

func makeCommissions() transfer.Commissions {
	return transfer.Commissions{
		FromInner: func(val int64) int64 {
			return int64(math.Max(float64(val*5/1000), 10_00))
		},
		ToInner: func(val int64) int64 {
			return 0
		},
		FromOuterToOuter: func(val int64) int64 {
			return int64(math.Max(float64(val*15/1000), 30_00))
		},
	}
}

func TestService_Card2Card(t *testing.T) {
	type fields struct {
		CardSvc        *card.Service
		TransactionSvc *transaction.Service
		commissions    transfer.Commissions
	}
	type args struct {
		from   string
		to     string
		amount int64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantTotal int64
		wantError error
	}{
		{
			name: "Transfer negative sum",
			fields: fields{
				CardSvc:        nil,
				TransactionSvc: nil,
				commissions:    transfer.Commissions{},
			},
			args: args{
				from:   "5106 2107 0000 0000",
				to:     "5106 2105 0000 0002",
				amount: -500_00,
			},
			wantTotal: 0,
			wantError: transfer.NonPositiveAmount,
		},
		{
			name: "Transfer zero sum",
			fields: fields{
				CardSvc:        nil,
				TransactionSvc: nil,
				commissions:    transfer.Commissions{},
			},
			args: args{
				from:   "5106 2107 0000 0000",
				to:     "5106 2105 0000 0002",
				amount: 0,
			},
			wantTotal: 0,
			wantError: transfer.NonPositiveAmount,
		},
		{
			name: "Inner success",
			fields: fields{
				CardSvc:        makeCardService2Cards(),
				TransactionSvc: transaction.NewService(),
				commissions:    makeCommissions(),
			},
			args: args{
				from:   "5106 2107 0000 0000",
				to:     "5106 2105 0000 0002",
				amount: 500_00,
			},
			wantTotal: 510_00,
			wantError: nil,
		},
		{
			name: "Inner not enough",
			fields: fields{
				CardSvc:        makeCardService2Cards(),
				TransactionSvc: transaction.NewService(),
				commissions:    makeCommissions(),
			},
			args: args{
				from:   "5106 2107 0000 0000",
				to:     "5106 2105 0000 0002",
				amount: 1000_00,
			},
			wantTotal: 1010_00,
			wantError: transfer.NotEnoughMoney,
		},
		{
			name: "Inner-outer success",
			fields: fields{
				CardSvc:        makeCardService1Card(),
				TransactionSvc: transaction.NewService(),
				commissions:    makeCommissions(),
			},
			args: args{
				from:   "5106 2107 0000 0000",
				to:     "2105 0000 5106 0002",
				amount: 500_00,
			},
			wantTotal: 510_00,
			wantError: nil,
		},
		{
			name: "Inner-outer not enough",
			fields: fields{
				CardSvc:        makeCardService1Card(),
				TransactionSvc: transaction.NewService(),
				commissions:    makeCommissions(),
			},
			args: args{
				from:   "5106 2107 0000 0000",
				to:     "2107 0000 5106 0000",
				amount: 1000_00,
			},
			wantTotal: 1010_00,
			wantError: transfer.NotEnoughMoney,
		},
		{
			name: "Outer-inner success",
			fields: fields{
				CardSvc:        makeCardService1Card(),
				TransactionSvc: transaction.NewService(),
				commissions:    makeCommissions(),
			},
			args: args{
				from:   "2107 0000 5106 0000",
				to:     "5106 2107 0000 0000",
				amount: 1000_00,
			},
			wantTotal: 1000_00,
			wantError: nil,
		},
		{
			name: "Outer success",
			fields: fields{
				CardSvc:        makeEmptyCardService(),
				TransactionSvc: transaction.NewService(),
				commissions:    makeCommissions(),
			},
			args: args{
				from:   "2107 0000 5106 0000",
				to:     "2107 5106 0000 0000",
				amount: 1000_00,
			},
			wantTotal: 1030_00,
			wantError: nil,
		},
		{
			name: "From card not found",
			fields: fields{
				CardSvc:        makeEmptyCardService(),
				TransactionSvc: transaction.NewService(),
				commissions:    makeCommissions(),
			},
			args: args{
				from:   "5106 2107 0000 0000",
				to:     "2107 5106 0000 0000",
				amount: 1000_00,
			},
			wantTotal: 0,
			wantError: transfer.CardNotFound,
		},
		{
			name: "To card not found",
			fields: fields{
				CardSvc:        makeEmptyCardService(),
				TransactionSvc: transaction.NewService(),
				commissions:    makeCommissions(),
			},
			args: args{
				from:   "5106 2107 0000 0000",
				to:     "5106 2105 0000 0002",
				amount: 1000_00,
			},
			wantTotal: 0,
			wantError: transfer.CardNotFound,
		},
	}

	for _, tt := range tests {
		s := transfer.NewService(tt.fields.CardSvc, tt.fields.TransactionSvc, tt.fields.commissions)
		gotTotal, gotError := s.Card2Card(tt.args.from, tt.args.to, tt.args.amount)
		if gotTotal != tt.wantTotal {
			t.Errorf("%v: gotTotal = %v, want %v", tt.name, gotTotal, tt.wantTotal)
		}
		if gotError != tt.wantError {
			t.Errorf("%v: gotError = %v, want %v", tt.name, gotError, tt.wantError)
		}
	}
}
