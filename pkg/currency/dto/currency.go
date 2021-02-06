package dto

import "math"

type RateXmlDTO struct {
	XMLName  string `xml:"Valute"`
	NumCode  string
	CharCode string
	Nominal  int64
	Name     string
	Value    float64
}

type RateListXmlDTO struct {
	XMLName string       `xml:"ValCurs"`
	Rates   []RateXmlDTO `xml:"Valute"`
}

func (d RateXmlDTO) ValueInCents() int64 {
	return int64(math.Round(d.Value*100)) / d.Nominal
}

type RateJsonDTO struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Value int64  `json:"value"`
}
