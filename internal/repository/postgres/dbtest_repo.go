package postgres

import (
	"context"
	"database/sql"
)

type DBTestRepo struct {
	db *sql.DB
}

func (r *DBTestRepo) InsertLine(ctx context.Context, line string) error {
	_, err := r.db.ExecContext(ctx, `insert into dbtest_logs(line) values ($1)`, line)
	return err
}

