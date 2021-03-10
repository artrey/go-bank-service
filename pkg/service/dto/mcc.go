package dto

import (
	"github.com/artrey/go-bank-service/pkg/models"
)

type Mcc struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

func FromModelMcc(m models.Mcc) *Mcc {
	return &Mcc{
		Id:   m.Id,
		Text: m.Text,
	}
}
