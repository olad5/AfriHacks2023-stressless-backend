package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra"
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

func (m *MongoMetricRepository) GetMetricById(ctx context.Context, metricId primitive.ObjectID) (domain.Metric, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoMetric := mongoMetric{}

	filter := bson.M{
		"_id": metricId,
	}
	err := m.metrics.FindOne(ctx, filter).Decode(&mongoMetric)
	if err != nil {
		m.logger.Error("failed to find metric by id: %w", zap.Error(err))
		return domain.Metric{}, infra.ErrMetricNotFound
	}
	return toDomainMetric(mongoMetric), nil
}

func (m *MongoMetricRepository) GetUserTodayLogIfExists(ctx context.Context, userId primitive.ObjectID) (domain.Metric, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoMetric := mongoMetric{}
	startTime, endTime := getDayBounds()

	filter := bson.M{
		"owner_id": userId,
		"created_at": bson.M{
			"$gte": primitive.NewDateTimeFromTime(startTime),
			"$lt":  primitive.NewDateTimeFromTime(endTime),
		},
	}

	err := m.metrics.FindOne(ctx, filter).Decode(&mongoMetric)
	if err != nil {
		m.logger.Error("failed to find mongo metric for today: %w", zap.Error(err))
		return domain.Metric{}, infra.ErrMetricNotFound
	}
	return toDomainMetric(mongoMetric), nil
}

func getDayBounds() (time.Time, time.Time) {
	now := time.Now()
	year, month, day := now.Date()

	earliest := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	latest := time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())

	return earliest, latest
}

func (m *MongoMetricRepository) GetRecentMetricsByUserId(ctx context.Context, userId primitive.ObjectID) ([]domain.Metric, error) {
	// TODO:TODO: I think this method has issues
	mongoMetrics := []*mongoMetric{}
	filter := bson.M{"owner_id": userId}
	cursor, err := m.metrics.Find(ctx, filter)
	if err != nil {
		m.logger.Error("failed retrieve recent metrics by user id: %w", zap.Error(err))
		return []domain.Metric{}, err
	}
	for cursor.Next(ctx) {
		var mm mongoMetric
		err := cursor.Decode(&mm)
		if err != nil {
			m.logger.Error("failed to decode metric in list of metrics : %w", zap.Error(err))
			return []domain.Metric{}, err
		}
		mongoMetrics = append(mongoMetrics, &mm)
	}
	if err := cursor.Err(); err != nil {
		return []domain.Metric{}, err
	}
	cursor.Close(ctx)

	result := []domain.Metric{}
	for _, element := range mongoMetrics {
		result = append(result, toDomainMetric(*element))
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
