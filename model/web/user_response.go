package web

import "time"

type UserResponse struct {
	Id        int
	Username  string
	CreatedAt time.Time
}
