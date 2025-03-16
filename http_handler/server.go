package http_handler

import (
	"fmt"

	"github.com/yuriyfomin17/bad-code-review/service"
)

type HttpServer struct {
	orderService service.OrderService
	userService  service.UserService
}

func NewHttpServer(httpTimeout, numWorkers, batchOrderNumber int) (*HttpServer, error) {
	us, err := service.NewUserService(httpTimeout, numWorkers)
	if err != nil {
		return nil, fmt.Errorf("error creating user service: %w", err)
	}
	orderService, err := service.NewOrderService(us, batchOrderNumber)
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
