package dto

import "math"

type RateDTO struct {
	XMLName  string `xml:"Valute"`
	NumCode  string
	CharCode string
	Nominal  int64
	Name     string
	Value    float64
}

type RateListDTO struct {
	XMLName string    `xml:"ValCurs"`
	Rates   []RateDTO `xml:"Valute"`
}

func (d RateDTO) ValueInCents() int64 {
	return int64(math.Round(d.Value * 100)) / d.Nominal
}
