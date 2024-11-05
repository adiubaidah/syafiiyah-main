package persistence

import (
	"context"
	"fmt"
)

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
