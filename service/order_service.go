package service

import (
	"bad-code-review/model"
)

type OrderServiceImpl struct {
	numOrdersIdsToProcess int
}

func NewOrderService(batchOrderNumber int) (OrderService, error) {

	return OrderServiceImpl{
		numOrdersIdsToProcess: batchOrderNumber,
	}, nil
}

func (os OrderServiceImpl) TransformAndSplitOrders(orderIDs []string) [][]model.Order {
	var batches [][]model.Order
	batch := make([]model.Order, 0, os.numOrdersIdsToProcess)
	ct := 0
	for _, id := range orderIDs {
		if ct < os.numOrdersIdsToProcess {
			batch = append(batch, model.Order{ID: id, UserID: "user" + id})
		} else if ct == os.numOrdersIdsToProcess {
			batches = append(batches, batch)
			batch = make([]model.Order, 0, os.numOrdersIdsToProcess)
		}
		ct += 1
		ct %= os.numOrdersIdsToProcess
	}
	if len(batch) > 0 {
		batches = append(batches, batch)
	}
	return batches
}
