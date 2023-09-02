package utils

// func WithTx(ctx context.Context, opts *sql.TxOptions, q mysql.Querier, db *sql.DB, f func(mysql.Querier) error) error {
// 	w, ok := q.(interface {
// 		WithTx(tx *sql.Tx) *mysql.Queries
// 	})
// 	var tx *sql.Tx
// 	if ok {
// 		var err error
// 		tx, err = db.BeginTx(ctx, opts)
// 		if err != nil {
// 			return fmt.Errorf("WithTx: %w", err)
// 		}
// 		defer tx.Rollback()
// 		q = w.WithTx(tx)
// 	}
// 	err := f(q)
// 	if err != nil {
// 		return fmt.Errorf("WithTx: %w", err)
// 	}
// 	if tx != nil {
// 		err := tx.Commit()
// 		if err != nil {
// 			return fmt.Errorf("WithTx: %w", err)
// 		}
// 	}
// 	return nil
// }
