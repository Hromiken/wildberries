package pgdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"order-notification/internal/entity"
	"order-notification/pkg/postgres"
)

type OrdersRepo struct {
	*postgres.Postgres
}

func NewOrdersRepo(pg *postgres.Postgres) *OrdersRepo {
	return &OrdersRepo{pg}
}

func (r *OrdersRepo) GetOrderByUID(ctx context.Context, uid uuid.UUID) (entity.Order, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("public.order").
		Where("order_uid = ?", uid).
		ToSql()
	if err != nil {
		return entity.Order{}, fmt.Errorf("failed to build order query: %w", err)
	}

	var (
		order                    entity.Order
		payment, delivery, items []byte
	)

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&delivery,
		&payment,
		&items,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, fmt.Errorf("order not found")
		}
		return entity.Order{}, fmt.Errorf("failed to scan order: %w", err)
	}

	err = json.Unmarshal(payment, &order.Payment)
	if err != nil {
		return entity.Order{}, fmt.Errorf("failed to unmarshal payment: %w", err)
	}

	err = json.Unmarshal(items, &order.Items)
	if err != nil {
		return entity.Order{}, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	err = json.Unmarshal(delivery, &order.Delivery)
	if err != nil {
		return entity.Order{}, fmt.Errorf("failed to unmarshal delivery: %w", err)
	}

	return order, nil
}

func (r *OrdersRepo) SaveOrder(ctx context.Context, order entity.Order) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	payment, err := json.Marshal(order.Payment)
	if err != nil {
		return fmt.Errorf("failed to marshal payment: %w", err)
	}
	delivery, err := json.Marshal(order.Delivery)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery: %w", err)
	}
	items, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	sql, args, err := r.Builder.
		Insert("public.order").
		Columns(
			"order_uid", "track_number", "entry",
			"delivery", "payment", "items",
			"locale", "internal_signature", "customer_id",
			"delivery_service", "shardkey", "sm_id",
			"date_created", "oof_shard",
		).
		Values(
			order.OrderUID, order.TrackNumber, order.Entry,
			delivery, payment, items,
			order.Locale, order.InternalSignature, order.CustomerID,
			order.DeliveryService, order.ShardKey, order.SmID,
			order.DateCreated, order.OofShard,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build order insert query: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("order with uid %s already exists: %w", order.OrderUID, err)
		}
		return fmt.Errorf("failed to insert order: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
