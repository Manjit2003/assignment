package model

import "time"

type TodoItem struct {
	ID          string    `json:"id"`
	UserID      string    `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type TodoPatch struct {
	Title       *string
	Description *string
	Status      *string
}
