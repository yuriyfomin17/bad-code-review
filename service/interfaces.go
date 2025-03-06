package service

import (
	"bad-code-review/model"
	"context"
)

type UserService interface {
	FetchUserDetailsBatch(ctx context.Context, orderIDs []string) ([]model.User, error)
}

type OrderService interface {
	TransformAndSplitOrders(orderIDs []string) [][]model.Order
}
