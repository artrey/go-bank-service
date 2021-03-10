package dto

import (
	"github.com/artrey/go-bank-service/pkg/models"
)

type Icon struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Uri   string `json:"uri"`
}

func FromModelIcon(i models.Icon) *Icon {
	return &Icon{
		Id:    i.Id,
		Title: i.Title,
		Uri:   i.Uri,
	}
}
