package infra

import (
	"context"
	"errors"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByUserId(ctx context.Context, userId primitive.ObjectID) (domain.User, error)
}
