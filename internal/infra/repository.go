package infra

import (
	"context"
	"errors"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrMetricNotFound         = errors.New("metric not found")
	ErrRecommendationNotFound = errors.New("recommendation not found")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByUserId(ctx context.Context, userId primitive.ObjectID) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
	UpdateUserLastMetricLog(ctx context.Context, user domain.User) error
}

type MetricRepository interface {
	CreateMetric(ctx context.Context, metric domain.Metric) error
	GetUserTodayLogIfExists(ctx context.Context, userId primitive.ObjectID) (domain.Metric, error)
	UpdateMetricById(ctx context.Context, metric domain.Metric) error
	GetMetricById(ctx context.Context, metricId primitive.ObjectID) (domain.Metric, error)
	GetRecentMetricsByUserId(ctx context.Context, userId primitive.ObjectID) ([]domain.Metric, error)
}

type RecommendationRepository interface {
	CreateRecommendation(ctx context.Context, metric domain.Recommendation) error
	UpdateRecommendationById(ctx context.Context, metric domain.Recommendation) error
	GetRecommendationById(ctx context.Context, metricId primitive.ObjectID) (domain.Recommendation, error)
	GetRecommendationByMetricId(ctx context.Context, metricId primitive.ObjectID, metricType string) (domain.Recommendation, error)
}
