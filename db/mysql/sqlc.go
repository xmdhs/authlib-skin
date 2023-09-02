//go:generate sqlc generate
package mysql

import (
	"context"
	"database/sql"
	"fmt"
)

type QuerierWithTx interface {
	Querier
	Tx(ctx context.Context, f func(Querier) error) error
}

var _ QuerierWithTx = (*Queries)(nil)

func (q *Queries) Tx(ctx context.Context, f func(Querier) error) error {
	db, ok := q.db.(*sql.DB)
	if !ok {
		return fmt.Errorf("not *sql.DB")
	}
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	nq := q.WithTx(tx)
	err = f(nq)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
