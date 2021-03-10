package dto

import (
	"github.com/artrey/go-bank-service/pkg/models"
)

type RequestCards struct {
	ClientId int64 `json:"clientId"`
}

type Card struct {
	Id        int64  `json:"id"`
	Number    string `json:"number"`
	Balance   int64  `json:"balance"`
	Issuer    string `json:"issuer"`
	Holder    string `json:"holder"`
	OwnerId   int64  `json:"ownerId"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"createdAt"`
}

func FromModelCard(c models.Card) *Card {
	return &Card{
		Id:        c.Id,
		Number:    c.Number,
		Balance:   c.Balance,
		Issuer:    c.Issuer,
		Holder:    c.Holder,
		OwnerId:   c.OwnerId,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
	}
}
