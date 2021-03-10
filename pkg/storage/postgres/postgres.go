package postgres

import (
	"context"
	"github.com/artrey/go-bank-service/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
	ctx context.Context
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	s := new(Storage)
	db, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	s.db = db
	s.ctx = ctx
	return s, nil
}

func (s *Storage) GetCardsByClientId(id int64) ([]models.Card, error) {
	sql := `
SELECT id, number, balance, issuer, holder, owner_id, status, EXTRACT(EPOCH FROM created_at)::bigint
FROM cards
WHERE owner_id = $1
`
	rows, err := s.db.Query(s.ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]models.Card, 0)
	for rows.Next() {
		var c models.Card
		err = rows.Scan(&c.Id, &c.Number, &c.Balance, &c.Issuer, &c.Holder, &c.OwnerId, &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (s *Storage) GetTransactionsByCardId(id int64) ([]models.Transaction, error) {
	sql := `
SELECT id, from_id, to_id, sum, mcc_id, icon_id, description, EXTRACT(EPOCH FROM created_at)::bigint
FROM transactions
WHERE from_id = $1 or to_id = $1
`
	rows, err := s.db.Query(s.ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		var t models.Transaction
		err = rows.Scan(&t.Id, &t.FromId, &t.ToId, &t.Sum, &t.MccId, &t.IconId, &t.Description, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
