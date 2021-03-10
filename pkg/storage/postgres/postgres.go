package postgres

import (
	"context"
	"github.com/artrey/go-bank-service/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db  *pgxpool.Pool
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

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) GetCardsByClientId(id int64) ([]models.Card, error) {
	sql := `
SELECT id, number, balance, issuer, holder, owner_id, status, EXTRACT(EPOCH FROM created_at)::bigint
FROM cards
WHERE owner_id = $1
LIMIT 50
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

func (s *Storage) GetCardById(id int64) (models.Card, error) {
	sql := `
SELECT id, number, balance, issuer, holder, owner_id, status, EXTRACT(EPOCH FROM created_at)::bigint
FROM cards
WHERE id = $1
LIMIT 1
`
	var c models.Card
	err := s.db.QueryRow(s.ctx, sql, id).Scan(
		&c.Id, &c.Number, &c.Balance, &c.Issuer, &c.Holder, &c.OwnerId, &c.Status, &c.CreatedAt)
	if err != nil {
		return models.Card{}, err
	}
	return c, nil
}

func (s *Storage) GetTransactionsByCardId(id int64) ([]models.Transaction, error) {
	sql := `
SELECT id, from_id, to_id, sum, mcc_id, icon_id, description, EXTRACT(EPOCH FROM created_at)::bigint
FROM transactions
WHERE from_id = $1 or to_id = $1
LIMIT 50
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

func (s *Storage) GetIconById(id int64) (models.Icon, error) {
	sql := `
SELECT id, title, uri
FROM icons
WHERE id = $1
LIMIT 1
`
	var icon models.Icon
	err := s.db.QueryRow(s.ctx, sql, id).Scan(&icon.Id, &icon.Title, &icon.Uri)
	if err != nil {
		return models.Icon{}, err
	}
	return icon, nil
}

func (s *Storage) GetMccById(id string) (models.Mcc, error) {
	sql := `
SELECT id, text
FROM mccs
WHERE id = $1
LIMIT 1
`
	var mcc models.Mcc
	err := s.db.QueryRow(s.ctx, sql, id).Scan(&mcc.Id, &mcc.Text)
	if err != nil {
		return models.Mcc{}, err
	}
	return mcc, nil
}

func (s *Storage) GetMostPopularSpendingByCard(id int64) (models.MostPopularSpending, error) {
	sql := `
SELECT COUNT(t.id) as count, t.description, i.uri
FROM transactions as t
INNER JOIN icons as i on i.id = t.icon_id
WHERE from_id = $1
GROUP BY t.description, t.icon_id, i.uri
ORDER BY count DESC
LIMIT 1
`
	var spent models.MostPopularSpending
	err := s.db.QueryRow(s.ctx, sql, id).Scan(&spent.Count, &spent.Description, &spent.IconUri)
	if err != nil {
		return models.MostPopularSpending{}, err
	}
	return spent, nil
}

func (s *Storage) GetMostExpensiveSpendingByCard(id int64) (models.MostExpensiveSpending, error) {
	sql := `
SELECT SUM(t.sum) as sum, t.description, i.uri
FROM transactions as t
INNER JOIN icons as i on i.id = t.icon_id
WHERE from_id = $1
GROUP BY t.description, t.icon_id, i.uri
ORDER BY sum
LIMIT 1
`
	var spent models.MostExpensiveSpending
	err := s.db.QueryRow(s.ctx, sql, id).Scan(&spent.Sum, &spent.Description, &spent.IconUri)
	if err != nil {
		return models.MostExpensiveSpending{}, err
	}
	spent.Sum = -spent.Sum
	return spent, nil
}
