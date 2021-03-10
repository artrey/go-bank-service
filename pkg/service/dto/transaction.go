package dto

import (
	"github.com/artrey/go-bank-service/pkg/models"
)

type Transaction struct {
	Id          int64   `json:"id"`
	From        *Card   `json:"from"`
	To          *Card   `json:"to"`
	Sum         int64   `json:"sum"`
	Mcc         *Mcc    `json:"mcc"`
	Icon        *Icon   `json:"icon"`
	Description *string `json:"description"`
	CreatedAt   int64   `json:"createdAt"`
}

func FromModelTransaction(t models.Transaction) *Transaction {
	return &Transaction{
		Id:          t.Id,
		Sum:         t.Sum,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
	}
}

type TransactionBuilder struct {
	t *Transaction
}

func NewTransactionBuilder(t *Transaction) *TransactionBuilder {
	return &TransactionBuilder{
		t: t,
	}
}

func (b *TransactionBuilder) Build() *Transaction {
	return b.t
}

func (b *TransactionBuilder) SetFrom(c *Card) *TransactionBuilder {
	b.t.From = c
	return b
}

func (b *TransactionBuilder) SetTo(c *Card) *TransactionBuilder {
	b.t.To = c
	return b
}

func (b *TransactionBuilder) SetMcc(m *Mcc) *TransactionBuilder {
	b.t.Mcc = m
	return b
}

func (b *TransactionBuilder) SetIcon(i *Icon) *TransactionBuilder {
	b.t.Icon = i
	return b
}
