package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID
	Email     string
	FirstName string
	LastName  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
