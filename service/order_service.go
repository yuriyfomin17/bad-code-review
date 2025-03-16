package service

import (
	"bad-code-review/model"
	"context"
	"log"
)

type OrderServiceImpl struct {
	userService           UserService
	numOrdersIdsToProcess int
}

func NewOrderService(us UserService, batchOrderNumber int) (OrderService, error) {

	return OrderServiceImpl{
		userService:           us,
		numOrdersIdsToProcess: batchOrderNumber,
	}, nil
}

func (os OrderServiceImpl) TransformOrderIdsToOrders(ctx context.Context, orderIDs []string) []model.Order {
	batchedOrderIds := os.transformToBatchOrderIds(orderIDs)
	userDetailsMap := os.fetchUserDetails(ctx, orderIDs, batchedOrderIds)

	orders := make([]model.Order, 0, len(orderIDs))
	for _, orderId := range orderIDs {
		userId := "user" + orderId
		userInMap, ok := userDetailsMap[userId]
		if !ok {
			continue
		}
		orders = append(orders, model.Order{ID: orderId, UserID: userId, User: userInMap})
	}

	return orders
}

func (os OrderServiceImpl) transformToBatchOrderIds(orderIDs []string) [][]string {
	var batches [][]string
	batch := make([]string, 0, os.numOrdersIdsToProcess)
	ct := 0
	for _, id := range orderIDs {
		if ct < os.numOrdersIdsToProcess {
			batch = append(batch, id)
		} else if ct == os.numOrdersIdsToProcess {
			batches = append(batches, batch)
			batch = make([]string, 0, os.numOrdersIdsToProcess)
		}
		ct += 1
		ct %= os.numOrdersIdsToProcess
	}
	if len(batch) > 0 {
		batches = append(batches, batch)
	}
	return batches
}

func (os OrderServiceImpl) fetchUserDetails(ctx context.Context, orderIDs []string, batchedOrderIds [][]string) map[string]model.User {
	userDetailsMap := make(map[string]model.User, len(orderIDs))
	for _, batchOrderIds := range batchedOrderIds {
		orderIdsToProcess := make([]string, 0, len(batchOrderIds))
		for _, orderId := range batchOrderIds {
			orderIdsToProcess = append(orderIdsToProcess, orderId)
		}
		usersBatch, err := os.userService.FetchUserDetailsBatch(ctx, orderIdsToProcess)
		if err != nil {
			log.Printf("error fetching user details for batchOrderIds [IDs: %v]: %v", orderIdsToProcess, err)
			continue
		}
		for _, user := range usersBatch {
			userDetailsMap[user.ID] = user
		}
	}
	return userDetailsMap
}
