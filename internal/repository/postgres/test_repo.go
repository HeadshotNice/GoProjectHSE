package postgres

import (
	"context"
	"database/sql"
)

type TestRepo struct {
	db *sql.DB
}

// Hello intentionally touches repository layer to satisfy LR1 requirement that
// request processing goes through all 3 layers.
func (r *TestRepo) Hello(ctx context.Context) (string, error) {
	var one int
	if err := r.db.QueryRowContext(ctx, "select 1").Scan(&one); err != nil {
		return "", err
	}
	return "Hello!", nil
}

