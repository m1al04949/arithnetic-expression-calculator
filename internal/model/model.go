package model

import "time"

type ExpressionTab struct {
	ID         int
	Added      time.Time
	Expression string
	Status     string
	Resul      float64
}
