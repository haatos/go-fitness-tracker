package model

type Exercise struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	UserID *string `json:"user_id"`
}
