package cache_test

import (
	"github.com/google/uuid"
	"order-notification/internal/entity"
	"order-notification/internal/repo/cache"
	"testing"
)

func TestCache(t *testing.T) {
	c := cache.NewCache(10, 60)

	id := uuid.New()

	order := entity.Order{
		OrderUID: id, // ТУТ ОШИБКА
	}

	if got := c.Get(order.OrderUID); got != nil {
		t.Errorf("expected nil, got %+v", got)
	}

	c.Set(order.OrderUID, &order)

	got := c.Get(order.OrderUID)
	if got == nil {
		t.Fatalf("expected order, got nil")
	}
	if got.OrderUID != order.OrderUID {
		t.Errorf("expected %s, got %s", order.OrderUID, got.OrderUID)
	}
}
