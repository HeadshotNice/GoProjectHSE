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
		`insert into orders(user_id, status) values ($1, $2) returning id`,
		userID,
		entity.OrderStatusCreated,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *OrdersRepo) ListByUser(ctx context.Context, userID int64, activeOnly bool) ([]entity.Order, error) {
	query := `select id, user_id, status, created_at, updated_at, completed_at from orders where user_id = $1 order by id desc`
	if activeOnly {
		query = `select id, user_id, status, created_at, updated_at, completed_at
			from orders
			where user_id = $1 and status in ('created', 'packing', 'arriving')
			order by id desc`
	}

	rows, err := r.db.QueryContext(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.Order
	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.CreatedAt, &o.UpdatedAt, &o.CompletedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *OrdersRepo) UpdateStatus(ctx context.Context, orderID int64, status string) error {
	var completedExpr string
	if status == entity.OrderStatusCompleted || status == entity.OrderStatusCanceled {
		completedExpr = "now()"
	} else {
		completedExpr = "completed_at"
	}

	_, err := r.db.ExecContext(
		ctx,
		`update orders set status = $1, updated_at = now(), completed_at = `+completedExpr+` where id = $2`,
		status,
		orderID,
	)
	return err
}
