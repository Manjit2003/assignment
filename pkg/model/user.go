package model

import "time"

type User struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"` // omit from the responses
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
}
