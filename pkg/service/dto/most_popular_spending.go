package dto

import (
	"github.com/artrey/go-bank-service/pkg/models"
)

type MostPopularSpending struct {
	Description string `json:"description"`
	Count       int64  `json:"count"`
	IconUri     string `json:"iconUri"`
}

func FromModelMostPopularSpending(s models.MostPopularSpending) *MostPopularSpending {
	return &MostPopularSpending{
		Description: s.Description,
		Count:       s.Count,
		IconUri:     s.IconUri,
	}
}
