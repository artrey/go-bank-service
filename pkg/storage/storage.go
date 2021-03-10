package storage

import "github.com/artrey/go-bank-service/pkg/models"

type Interface interface {
	GetCardsByClientId(int64) ([]models.Card, error)
	GetTransactionsByCardId(int64) ([]models.Transaction, error)
}
