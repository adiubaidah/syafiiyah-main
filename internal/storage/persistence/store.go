package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	ListSantri(ctx context.Context, arg ListSantriParams) ([]ListSantriRow, error)
	ListUsers(ctx context.Context, arg ListUserParams) ([]ListUserRow, error)
	ListParents(ctx context.Context, arg ListParentParams) ([]ListParentRow, error)
	ListEmployees(ctx context.Context, arg ListEmployeesParams) ([]ListEmployeesRow, error)
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

func (q *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := q.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	qs := q.WithTx(tx)
	if err := fn(qs); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
