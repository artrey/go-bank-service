package service

import (
	"encoding/json"
	"errors"
	"github.com/artrey/go-bank-service/pkg/service/dto"
	"github.com/artrey/go-bank-service/pkg/storage"
	"log"
	"net/http"
)

type Server struct {
	storage storage.Interface
	mux     *http.ServeMux
}

func NewServer(storage storage.Interface, mux *http.ServeMux) *Server {
	return &Server{
		storage: storage,
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
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if r.Method != "GET" {
		err := errors.New("not allowed, use GET")
		log.Println(err)
		w.WriteHeader(405)
		_ = encoder.Encode(dto.MakeError("not-allowed", err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var requestCards dto.RequestCards
	err := decoder.Decode(&requestCards)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		_ = encoder.Encode(dto.MakeError("invalid-data", err))
		return
	}

	cards, err := s.storage.GetCardsByClientId(requestCards.ClientId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		_ = encoder.Encode(dto.MakeUnknownError(err))
		return
	}

	dtos := make([]*dto.Card, len(cards))
	for i, card := range cards {
		dtos[i] = dto.FromModelCard(card)
	}
	err = encoder.Encode(dtos)

	if err != nil {
		w.WriteHeader(500)
		_ = encoder.Encode(dto.MakeUnknownError(err))
	}
}

func (s *Server) getTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if r.Method != "GET" {
		err := errors.New("not allowed, use GET")
		log.Println(err)
		w.WriteHeader(405)
		_ = encoder.Encode(dto.MakeError("not-allowed", err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var requestTransactions dto.RequestTransactions
	err := decoder.Decode(&requestTransactions)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		_ = encoder.Encode(dto.MakeError("invalid-data", err))
		return
	}

	transactions, err := s.storage.GetTransactionsByCardId(requestTransactions.CardId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		_ = encoder.Encode(dto.MakeUnknownError(err))
		return
	}

	dtos := make([]*dto.Transaction, len(transactions))
	for i, t := range transactions {
		b := dto.NewTransactionBuilder(dto.FromModelTransaction(t))
		if t.FromId != nil {
			c, err := s.storage.GetCardById(*t.FromId)
			if err != nil {
				w.WriteHeader(500)
				_ = encoder.Encode(dto.MakeUnknownError(err))
			}
			b.SetFrom(dto.FromModelCard(c))
		}
		if t.ToId != nil {
			c, err := s.storage.GetCardById(*t.ToId)
			if err != nil {
				w.WriteHeader(500)
				_ = encoder.Encode(dto.MakeUnknownError(err))
			}
			b.SetTo(dto.FromModelCard(c))
		}
		if t.MccId != nil {
			m, err := s.storage.GetMccById(*t.MccId)
			if err != nil {
				w.WriteHeader(500)
				_ = encoder.Encode(dto.MakeUnknownError(err))
			}
			b.SetMcc(dto.FromModelMcc(m))
		}
		if t.IconId != nil {
			i, err := s.storage.GetIconById(*t.IconId)
			if err != nil {
				w.WriteHeader(500)
				_ = encoder.Encode(dto.MakeUnknownError(err))
			}
			b.SetIcon(dto.FromModelIcon(i))
		}
		dtos[i] = b.Build()
	}
	err = encoder.Encode(dtos)

	if err != nil {
		w.WriteHeader(500)
		_ = encoder.Encode(dto.MakeUnknownError(err))
	}
}
