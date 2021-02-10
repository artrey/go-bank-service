package dto

import "github.com/artrey/go-bank-service/pkg/card"

type Card struct {
	Id           int64    `json:"id"`
	CardHolderId int64    `json:"cardHolderId"`
	Type         card.Typ `json:"type"`
	Issuer       string   `json:"issuer"`
	Balance      int64    `json:"balance"`
	Currency     string   `json:"currency"`
	Number       string   `json:"number"`
	Icon         string   `json:"icon"`
}

func FromServiceCard(c *card.Card) *Card {
	return &Card{
		Id:           c.Id,
		CardHolderId: c.CardHolderId,
		Type:         c.Type,
		Issuer:       c.Issuer,
		Balance:      c.Balance,
		Currency:     c.Issuer,
		Number:       c.Number,
		Icon:         c.Icon,
	}
}

type AddCard struct {
	CardHolderId int64    `json:"cardHolderId"`
	Issuer       string   `json:"issuer"`
	Type         card.Typ `json:"type"`
}
