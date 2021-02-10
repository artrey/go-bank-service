package card

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

type Typ string

const (
	Plastic Typ = "plastic"
	Virtual Typ = "virtual"
)

type Card struct {
	Id           int64
	CardHolderId int64
	Type         Typ
	Issuer       string
	Balance      int64
	Currency     string
	Number       string
	Icon         string
}

func (c *Card) Withdraw(amount int64) bool {
	if c.Balance < amount {
		return false
	}
	c.Balance -= amount
	return true
}

func (c *Card) AddMoney(amount int64) {
	c.Balance += amount
}

type Service struct {
	BankName     string
	issuerNumber string
	cards        []*Card
	mu           sync.RWMutex
}

func NewService(bankName, issuerNumber string) *Service {
	return &Service{
		BankName:     bankName,
		issuerNumber: issuerNumber,
		cards:        []*Card{},
	}
}

func (s *Service) All(ctx context.Context, cardHolderId int64) []*Card {
	cards := make([]*Card, 0)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, c := range s.cards {
		if c.CardHolderId == cardHolderId {
			cards = append(cards, c)
		}
	}
	return cards
}

func (s *Service) CheckOwning(cardNumber string) bool {
	return strings.HasPrefix(cardNumber, s.issuerNumber)
}

func (s *Service) GenerateNumber() string {
	// TODO: make actual algorithm
	first := rand.Int() % 7
	second := 7 - first
	return fmt.Sprintf("%s0%d 0000 000%d", s.issuerNumber, first, second)
}

func (s *Service) FindCard(ctx context.Context, number string) *Card {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, c := range s.cards {
		if c.Number == number {
			return c
		}
	}
	return nil
}

func (s *Service) FindCardsByHolder(ctx context.Context, id int64) []*Card {
	cards := make([]*Card, 0)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, c := range s.cards {
		if c.CardHolderId == id {
			cards = append(cards, c)
		}
	}
	return cards
}

func (s *Service) Issue(ctx context.Context, issuer string, cardHolderId int64, typ Typ,
	startBalance int64, currency, number, icon string) *Card {
	var id int64 = 1

	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.cards) > 0 {
		id = s.cards[len(s.cards)-1].Id + 1
	}

	card := &Card{
		Id:           id,
		CardHolderId: cardHolderId,
		Type:         typ,
		Issuer:       issuer,
		Balance:      startBalance,
		Currency:     currency,
		Number:       number,
		Icon:         icon,
	}
	s.cards = append(s.cards, card)

	return card
}
