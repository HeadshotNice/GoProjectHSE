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
		`insert into documents(user_id, title, content, status) values ($1, $2, $3, $4) returning id`,
		userID,
		title,
		content,
		entity.DocumentStatusPendingReview,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *DocumentsRepo) ListByUser(ctx context.Context, userID int64) ([]entity.Document, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`select id, user_id, title, content, status, created_at, updated_at from documents where user_id = $1 order by id desc`,
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
			&document.UpdatedAt,
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

func (r *DocumentsRepo) UpdateStatus(ctx context.Context, documentID int64, status string) error {
	_, err := r.db.ExecContext(
		ctx,
		`update documents set status = $1, updated_at = now() where id = $2`,
		status,
		documentID,
	)
	return err
}
