package repo

import (
	"github.com/google/uuid"
	"order-notification/internal/entity"
)

type Cache interface {
	Get(key uuid.UUID) *entity.Order
	Set(key uuid.UUID, order entity.Order)
}
