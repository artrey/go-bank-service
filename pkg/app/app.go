package app

import (
	"encoding/json"
	"errors"
	"github.com/artrey/go-bank-service/pkg/app/dto"
	"github.com/artrey/go-bank-service/pkg/card"
	"github.com/artrey/go-bank-service/pkg/transaction"
	"log"
	"net/http"
)

type Server struct {
	cardSvc        *card.Service
	transactionSvc *transaction.Service
	mux            *http.ServeMux
}

func NewServer(cardSvc *card.Service, transactionSvc *transaction.Service, mux *http.ServeMux) *Server {
	return &Server{
		cardSvc: cardSvc,
		transactionSvc: transactionSvc,
		mux:     mux,
	}
}

func (s *Server) Init() {
	s.mux.HandleFunc("/cards", s.getCards)
	s.mux.HandleFunc("/transactions", s.getTransactions)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) getCards(w http.ResponseWriter, r *http.Request) {
	cards := s.cardSvc.All(r.Context(), 1)
	dtos := make([]*dto.Card, len(cards))
	for i, c := range cards {
		dtos[i] = dto.FromServiceCard(c)
	}

	w.Header().Add("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(dtos)

	if err != nil {
		_ = encoder.Encode(dto.Error{
			Code:    "unknown",
			Message: err.Error(),
		})
		w.WriteHeader(500)
	}
}

func (s *Server) getTransactions(w http.ResponseWriter, r *http.Request) {
	transactions := s.transactionSvc.All(r.Context(), 1)
	dtos := make([]*dto.Card, len(transactions))
	for i, c := range transactions {
		dtos[i] = dto.FromServiceCard(c)
	}

	w.Header().Add("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(dtos)

	if err != nil {
		_ = encoder.Encode(dto.Error{
			Code:    "unknown",
			Message: err.Error(),
		})
		w.WriteHeader(500)
	}
}

func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if r.Method != "POST" {
		err := errors.New("not allowed, use POST")
		log.Println(err)
		w.WriteHeader(405)
		_ = encoder.Encode(dto.Error{
			Code:    "not-allowed",
			Message: err.Error(),
		})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var dtoAddCard dto.AddCard
	err := decoder.Decode(&dtoAddCard)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		_ = encoder.Encode(dto.Error{
			Code:    "invalid-data",
			Message: err.Error(),
		})
		return
	}

	holderCards := s.cardSvc.FindCardsByHolder(r.Context(), dtoAddCard.CardHolderId)
	if len(holderCards) == 0 {
		err = errors.New("no cards")
		log.Println(err)
		w.WriteHeader(400)
		_ = encoder.Encode(dto.Error{
			Code:    "not-permitted",
			Message: err.Error(),
		})
		return
	}

	newCard := s.cardSvc.Issue(r.Context(), dtoAddCard.Issuer, dtoAddCard.CardHolderId,
		dtoAddCard.Type, 0, "RUB", s.cardSvc.GenerateNumber(), "https://...")

	w.WriteHeader(201)
	err = encoder.Encode(newCard)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		_ = encoder.Encode(dto.Error{
			Code:    "unknown",
			Message: err.Error(),
		})
	}
}

func (s *Server) editCard(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (s *Server) removeCard(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
