package domain

import "time"

type User struct {
	Id             int
	Username       string
	HashedPassword string
	CreatedAt      time.Time
}
