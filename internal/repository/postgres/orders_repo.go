package postgres

import (
	"context"
	"database/sql"

	"HSE/internal/entity"
)

type OrdersRepo struct {
	db *sql.DB
}

func (r *OrdersRepo) Create(ctx context.Context, userID int64) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(
		ctx,
		`insert into orders(user_id) values ($1) returning id`,
		userID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *OrdersRepo) ListByUser(ctx context.Context, userID int64) ([]entity.Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`select id, user_id, status, created_at from orders where user_id = $1 order by id desc`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.Order
	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

