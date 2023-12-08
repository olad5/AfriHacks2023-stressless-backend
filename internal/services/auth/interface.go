package auth

import (
	"context"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTClaims struct {
	ID    primitive.ObjectID
	Email string
}

type ctxKey int

const jwtKey ctxKey = 1

func SetJWTClaims(ctx context.Context, jwt JWTClaims) context.Context {
	return context.WithValue(ctx, jwtKey, jwt)
}

func GetJWTClaims(ctx context.Context) (JWTClaims, bool) {
	v, ok := ctx.Value(jwtKey).(JWTClaims)
	return v, ok
}

type AuthService interface {
	DecodeJWT(ctx context.Context, tokenString string) (JWTClaims, error)
	GenerateJWT(ctx context.Context, user domain.User) (string, error)
	IsUserLoggedIn(ctx context.Context, authHeader, userId string) bool
}
