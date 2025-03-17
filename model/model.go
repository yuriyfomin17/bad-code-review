package model

type Order struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	User   User   `json:"user"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
