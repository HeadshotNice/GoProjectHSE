package postgres

import (
	"context"
	"database/sql"

	"HSE/internal/entity"
)

type DocumentsRepo struct {
	db *sql.DB
}

func (r *DocumentsRepo) Create(ctx context.Context, userID int64, title, content string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(
		ctx,
		`insert into documents(user_id, title, content) values ($1, $2, $3) returning id`,
		userID,
		title,
		content,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *DocumentsRepo) ListByUser(ctx context.Context, userID int64) ([]entity.Document, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`select id, user_id, title, content, status, created_at from documents where user_id = $1 order by id desc`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.Document
	for rows.Next() {
		var document entity.Document
		if err := rows.Scan(
			&document.ID,
			&document.UserID,
			&document.Title,
			&document.Content,
			&document.Status,
			&document.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, document)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
