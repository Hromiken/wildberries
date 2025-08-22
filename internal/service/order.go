package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"order-notification/internal/entity"
	"order-notification/internal/repo/cache"
	"order-notification/internal/repo/pgdb"
	"order-notification/pkg/postgres"
)

type OrderService struct {
	Repo      *pgdb.OrdersRepo
	Cache     *cache.OrderCache
	validator *validator.Validate
}

func NewOrderService(pg *postgres.Postgres) *OrderService {
	return &OrderService{
		Repo:      pgdb.NewOrdersRepo(pg),
		Cache:     cache.NewCache(10, 10),
		validator: validator.New(),
	}
}

func (s *OrderService) SaveOrder(order entity.Order, ctx context.Context) error {

	if err := s.validator.Struct(order); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	if order.OrderUID == uuid.Nil {
		return errors.New("order_uid is empty")
	}
	if err := s.Repo.SaveOrder(context.Background(), order); err != nil {
		return err
	}
	s.Cache.Set(order.OrderUID, &order)
	return nil
}

func (s *OrderService) GetOrderByUID(ctx context.Context, uid uuid.UUID) (entity.Order, error) {
	if order := s.Cache.Get(uid); order != nil {
		return *order, nil
	}

	order, err := s.Repo.GetOrderByUID(ctx, uid)
	if err != nil {
		return entity.Order{}, err
	}

	s.Cache.Set(uid, &order)

	return order, nil
}
