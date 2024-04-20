package model

import "time"

type ExpressionTab struct {
	ID         int
	Added      time.Time
	Expression string
	Status     string
	Result     float64
}

type UsersTab struct {
	UserID    int
	Login     string
	Password  string
	CreatedAt time.Time
}
