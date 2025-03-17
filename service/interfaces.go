package service

import (
	"context"

	"github.com/yuriyfomin17/bad-code-review/model"
)

type UserService interface {
	FetchUserDetailsBatch(ctx context.Context, orderIDs []string) ([]model.User, error)
}

type OrderService interface {
	TransformOrderIdsToOrders(ctx context.Context, orderIDs []string) []model.Order
}
