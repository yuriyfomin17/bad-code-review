package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Order struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	User   User   `json:"user"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	http.HandleFunc("/orders", OrderHandler)
	http.ListenAndServe(":8080", nil)
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("ids")
	orderIDs := strings.Split(ids, ",")

	orders := FetchOrders(orderIDs)

	for i := range orders {
		user := FetchUserDetails(orders[i].UserID)
		orders[i].User = user
	}

	resp, _ := json.Marshal(orders)
	w.Write(resp)
}

func FetchOrders(orderIDs []string) []Order {
	var orders []Order
	for _, id := range orderIDs {
		orders = append(orders, Order{ID: id, UserID: "user" + id})
	}
	return orders
}

func FetchUserDetails(userID string) User {
	resp, _ := http.Get("http://user-service:8081/user?id=" + userID)
	body, _ := ioutil.ReadAll(resp.Body)
	var user User
	json.Unmarshal(body, &user)
	return user
}
