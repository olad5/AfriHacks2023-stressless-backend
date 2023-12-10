package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoMetricRepository struct {
	metrics *mongo.Collection
	logger  *zap.Logger
}

func NewMongoMetricRepo(ctx context.Context, mongoDatabase *mongo.Database, logger *zap.Logger) (*MongoMetricRepository, error) {
	metricsCollection := mongoDatabase.Collection("metrics")

	return &MongoMetricRepository{metrics: metricsCollection, logger: logger}, nil
}

func (m *MongoMetricRepository) CreateMetric(ctx context.Context, metric domain.Metric) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoMetric := toMongoMetric(metric)

	_, err := m.metrics.InsertOne(ctx, mongoMetric)
	if err != nil {
		m.logger.Error("failed to persist metric: %w", zap.Error(err))
		return fmt.Errorf("failed to persist metric: %w", err)
	}
	return nil
}

func (m *MongoMetricRepository) UpdateMetricById(ctx context.Context, metric domain.Metric) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoMetric := toMongoMetric(metric)

	filter := bson.M{"_id": metric.ID}
	updatedDoc := bson.M{
		"$set": mongoMetric,
	}
	_, err := m.metrics.UpdateOne(ctx, filter, updatedDoc)
	if err != nil {
		m.logger.Error("failed to update metric: %w", zap.Error(err))
		return fmt.Errorf("failed to update metric: %w", err)
	}
	return nil
}

func (m *MongoMetricRepository) GetRecentMetricsByUserId(ctx context.Context, userId primitive.ObjectID) ([]domain.Metric, error) {
	mongoMetrics := []mongoMetric{}
	err := m.metrics.FindOne(ctx, bson.M{"owner_id": userId}).Decode(&mongoMetrics)
	if err != nil {
		m.logger.Error("failed retrieve recent metrics by user id: %w", zap.Error(err))
		return []domain.Metric{}, err
	}
	result := []domain.Metric{}
	for _, element := range mongoMetrics {
		result = append(result, toDomainMetric(element))
	}

	return result, nil
}

type mongoMetric struct {
	ObjectID        primitive.ObjectID  `bson:"_id"`
	OwnerId         primitive.ObjectID  `bson:"owner_id"`
	StressLevel     int                 `bson:"stress_level"`
	Mood            domain.Mood         `bson:"mood"`
	SleepQuality    domain.SleepQuality `bson:"sleep_quality"`
	Feeling         string              `bson:"feeling"`
	StressLessScore int                 `bson:"stress_less_score"`
	CreatedAt       time.Time           `bson:"created_at"`
	UpdatedAt       time.Time           `bson:"updated_at"`
}

func toMongoMetric(metric domain.Metric) mongoMetric {
	return mongoMetric{
		ObjectID:        metric.ID,
		OwnerId:         metric.OwnerId,
		StressLevel:     metric.StressLevel,
		StressLessScore: metric.StressLessScore,
		SleepQuality:    metric.SleepQuality,
		Mood:            metric.Mood,
		Feeling:         metric.Feeling,
		CreatedAt:       metric.CreatedAt,
		UpdatedAt:       metric.UpdatedAt,
	}
}

func toDomainMetric(m mongoMetric) domain.Metric {
	return domain.Metric{
		ID:              m.ObjectID,
		OwnerId:         m.OwnerId,
		StressLevel:     m.StressLevel,
		StressLessScore: m.StressLessScore,
		SleepQuality:    m.SleepQuality,
		Mood:            m.Mood,
		Feeling:         m.Feeling,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
