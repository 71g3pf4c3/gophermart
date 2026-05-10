package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/71g3pf4c3/gophermart/pkg/postgres"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// OrderRepository implements order storage in PostgreSQL.
type OrderRepository struct {
	*postgres.Postgres
}

// NewOrderRepository creates a new OrderRepository.
func NewOrderRepository(pg *postgres.Postgres) *OrderRepository {
	return &OrderRepository{pg}
}

// CreateOrder inserts a new order. Returns ErrConflict if number already exists.
func (r *OrderRepository) CreateOrder(ctx context.Context, order entity.Order) error {
	sql, args, err := r.Builder.
		Insert("orders").
		Columns("number, user_id, status, uploaded_at").
		Values(order.Number, order.UserID, order.Status, order.UploadedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrderRepository - CreateOrder - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrConflict
		}
		return fmt.Errorf("OrderRepository - CreateOrder - Exec: %w", err)
	}

	return nil
}

// GetOrderByNumber returns a single order by its number.
func (r *OrderRepository) GetOrderByNumber(ctx context.Context, number string) (entity.Order, error) {
	sql, args, err := r.Builder.
		Select("number, user_id, status, accrual, uploaded_at").
		From("orders").
		Where(sq.Eq{"number": number}).
		ToSql()
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderRepository - GetOrderByNumber - Builder: %w", err)
	}

	var o entity.Order
	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&o.Number, &o.UserID, &o.Status, &o.Accrual, &o.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, entity.ErrOrderNotFound
		}
		return entity.Order{}, fmt.Errorf("OrderRepository - GetOrderByNumber - Scan: %w", err)
	}

	return o, nil
}

// GetOrdersByUser returns all orders for a user sorted by upload time desc.
func (r *OrderRepository) GetOrdersByUser(ctx context.Context, userID string) ([]entity.Order, error) {
	sql, args, err := r.Builder.
		Select("number, user_id, status, accrual, uploaded_at").
		From("orders").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("uploaded_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("OrderRepository - GetOrdersByUser - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("OrderRepository - GetOrdersByUser - Query: %w", err)
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		if err = rows.Scan(&o.Number, &o.UserID, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			return nil, fmt.Errorf("OrderRepository - GetOrdersByUser - Scan: %w", err)
		}
		orders = append(orders, o)
	}

	return orders, rows.Err()
}

// GetPendingOrders returns orders that still need accrual processing.
func (r *OrderRepository) GetPendingOrders(ctx context.Context) ([]entity.Order, error) {
	sql, args, err := r.Builder.
		Select("number, user_id, status, accrual, uploaded_at").
		From("orders").
		Where(sq.Or{sq.Eq{"status": entity.OrderStatusNew}, sq.Eq{"status": entity.OrderStatusProcessing}}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("OrderRepository - GetPendingOrders - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("OrderRepository - GetPendingOrders - Query: %w", err)
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		if err = rows.Scan(&o.Number, &o.UserID, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			return nil, fmt.Errorf("OrderRepository - GetPendingOrders - Scan: %w", err)
		}
		orders = append(orders, o)
	}

	return orders, rows.Err()
}

// UpdateOrderStatus updates status and accrual for an order.
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, number string, status entity.OrderStatus, accrual *float64) error {
	sql, args, err := r.Builder.
		Update("orders").
		Set("status", status).
		Set("accrual", accrual).
		Where(sq.Eq{"number": number}).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrderRepository - UpdateOrderStatus - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OrderRepository - UpdateOrderStatus - Exec: %w", err)
	}

	return nil
}
