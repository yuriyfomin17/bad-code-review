package http_handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *HttpServer) OrderHandler(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("ids")
	orderIDs := strings.Split(ids, ",")
	if len(orderIDs) == 0 {
		http.Error(w, "no orders found", http.StatusNotFound)
		return
	}
	orders := s.orderService.TransformOrderIdsToOrders(r.Context(), orderIDs)
	err := json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
