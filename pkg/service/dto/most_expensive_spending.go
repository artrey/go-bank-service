package dto

import (
	"github.com/artrey/go-bank-service/pkg/models"
)

type MostExpensiveSpending struct {
	Description string `json:"description"`
	Sum         int64  `json:"sum"`
	IconUri     string `json:"iconUri"`
}

func FromModelMostExpensiveSpending(s models.MostExpensiveSpending) *MostExpensiveSpending {
	return &MostExpensiveSpending{
		Description: s.Description,
		Sum:         s.Sum,
		IconUri:     s.IconUri,
	}
}
