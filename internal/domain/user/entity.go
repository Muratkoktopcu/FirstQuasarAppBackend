package user

import "time"

type User struct {
	ID        int64
	Email     string
	Password  string // hashed
	FullName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
