package http_handler

import (
	"bad-code-review/model"
	"bad-code-review/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpServer struct {
	orderService service.OrderService
	userService  service.UserService
}

func NewHttpServer(httpTimeout, numWorkers, batchOrderNumber int) (*HttpServer, error) {
	orderService, err := service.NewOrderService(batchOrderNumber)
	if err != nil {
		return nil, fmt.Errorf("error creating order service: %w", err)
	}
	userService, err := service.NewUserService(httpTimeout, numWorkers)
	if err != nil {
		return nil, fmt.Errorf("error creating user service: %w", err)
	}
	return &HttpServer{
		orderService: orderService,
		userService:  userService,
	}, nil
}

func (s *HttpServer) OrderHandler(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("ids")
	orderIDs := strings.Split(ids, ",")
	if len(orderIDs) == 0 {
		http.Error(w, "no orders found", http.StatusNotFound)
		return
	}

	var batchedOrders [][]model.Order
	batchedOrders = s.orderService.TransformAndSplitOrders(orderIDs)
	userDetailsMap := make(map[string]model.User, len(orderIDs))
	for _, batch := range batchedOrders {
		orderIdsToProcess := make([]string, 0, len(batch))
		for _, order := range batch {
			orderIdsToProcess = append(orderIdsToProcess, order.ID)
		}
		usersBatch, err := s.userService.FetchUserDetailsBatch(r.Context(), orderIdsToProcess)
		if err != nil {
			log.Printf("error fetching user details for batch [IDs: %v]: %v", orderIdsToProcess, err)
			continue
		}
		for _, user := range usersBatch {
			userDetailsMap[user.ID] = user
		}
	}
	var orders []model.Order
	for _, batch := range batchedOrders {
		orders = append(orders, batch...)
	}

	for _, order := range orders {
		userInMap, ok := userDetailsMap[order.UserID]
		if !ok {
			continue
		}
		order.User = userInMap
	}
	err := json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
