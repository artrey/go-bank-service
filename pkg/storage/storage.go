package storage

import "github.com/artrey/go-bank-service/pkg/models"

type Interface interface {
	GetCardsByClientId(int64) ([]models.Card, error)
	GetCardById(int64) (models.Card, error)
	GetTransactionsByCardId(int64) ([]models.Transaction, error)
	GetIconById(int64) (models.Icon, error)
	GetMccById(string) (models.Mcc, error)
	GetMostPopularSpendingByCard(int64) (models.MostPopularSpending, error)
	GetMostExpensiveSpendingByCard(int64) (models.MostExpensiveSpending, error)
}
