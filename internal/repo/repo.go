package repo

import (
	"context"
	"github.com/google/uuid"

	"order-notification/internal/entity"
)

type Repo interface {
	GetOrderByUID(ctx context.Context, uid uuid.UUID) (entity.Order, error)
}
