package postgres

import (
	"context"
	"database/sql"

	"HSE/internal/entity"
)

type UsersRepo struct {
	db *sql.DB
}

func (r *UsersRepo) Create(ctx context.Context, email, passwordHash string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(
		ctx,
		`insert into users(email, password_hash) values ($1, $2) returning id`,
		email,
		passwordHash,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UsersRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u entity.User
	err := r.db.QueryRowContext(
		ctx,
		`select id, email, password_hash, created_at from users where email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
