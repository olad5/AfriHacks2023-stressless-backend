package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MongoRecommendationRepository struct {
	recommendations *mongo.Collection
	logger          *zap.Logger
}

func NewMongoRecommendationRepo(ctx context.Context, mongoDatabase *mongo.Database, logger *zap.Logger) (*MongoRecommendationRepository, error) {
	recommendationsCollection := mongoDatabase.Collection("recommendations")

	return &MongoRecommendationRepository{recommendations: recommendationsCollection, logger: logger}, nil
}

func (m *MongoRecommendationRepository) CreateRecommendation(ctx context.Context, recommendation domain.Recommendation) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoRecommendation := toMongoRecommendation(recommendation)

	_, err := m.recommendations.InsertOne(ctx, mongoRecommendation)
	if err != nil {
		m.logger.Error("failed to persist recommendation: %w", zap.Error(err))
		return fmt.Errorf("failed to persist recommendation: %w", err)
	}
	return nil
}

func (m *MongoRecommendationRepository) UpdateRecommendationById(ctx context.Context, recommendation domain.Recommendation) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoRecommendation := toMongoRecommendation(recommendation)

	filter := bson.M{"_id": recommendation.ID}
	updatedDoc := bson.M{
		"$set": mongoRecommendation,
	}
	_, err := m.recommendations.UpdateOne(ctx, filter, updatedDoc)
	if err != nil {
		m.logger.Error("failed to update recommendation: %w", zap.Error(err))
		return fmt.Errorf("failed to update recommendation: %w", err)
	}
	return nil
}

func (m *MongoRecommendationRepository) GetRecommendationById(ctx context.Context, recommendationId primitive.ObjectID) (domain.Recommendation, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoRecommendation := mongoRecommendation{}

	filter := bson.M{
		"_id": recommendationId,
	}
	err := m.recommendations.FindOne(ctx, filter).Decode(&mongoRecommendation)
	if err != nil {
		m.logger.Error("failed to find recommendation by id: %w", zap.Error(err))
		return domain.Recommendation{}, infra.ErrRecommendationNotFound
	}
	return toDomainRecommendation(mongoRecommendation), nil
}

func (m *MongoRecommendationRepository) GetRecommendationByMetricId(ctx context.Context, metricId primitive.ObjectID, metricType string) (domain.Recommendation, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoRecommendation := mongoRecommendation{}

	filter := bson.M{
		"metric_id":   metricId,
		"metric_type": metricType,
	}
	err := m.recommendations.FindOne(ctx, filter).Decode(&mongoRecommendation)
	if err != nil {
		m.logger.Error("failed to find recommendation by metric id: %w", zap.Error(err))
		return domain.Recommendation{}, infra.ErrRecommendationNotFound
	}
	return toDomainRecommendation(mongoRecommendation), nil
}

type mongoRecommedationItem struct {
	Index    int    `bson:"index"`
	Heading  string `bson:"heading"`
	Text     string `bson:"text"`
	ImageUrl string `bson:"image_url"`
}

type mongoRecommendation struct {
	ObjectID   primitive.ObjectID       `bson:"_id"`
	MetricId   primitive.ObjectID       `bson:"metric_id"`
	MetricType string                   `bson:"metric_type"`
	Items      []mongoRecommedationItem `bson:"items"`
	CreatedAt  time.Time                `bson:"created_at"`
	UpdatedAt  time.Time                `bson:"updated_at"`
}

func toMongoRecommendation(recommendation domain.Recommendation) mongoRecommendation {
	var items []mongoRecommedationItem
	for _, item := range recommendation.Items {
		items = append(items, toMongoRecommendationItem(item))
	}
	return mongoRecommendation{
		ObjectID:   recommendation.ID,
		MetricId:   recommendation.MetricId,
		MetricType: recommendation.MetricType,
		Items:      items,
		CreatedAt:  recommendation.CreatedAt,
		UpdatedAt:  recommendation.UpdatedAt,
	}
}

func toMongoRecommendationItem(recommendationItem domain.RecommendationItem) mongoRecommedationItem {
	return mongoRecommedationItem{
		Text:     recommendationItem.Text,
		Index:    recommendationItem.Index,
		Heading:  recommendationItem.Heading,
		ImageUrl: recommendationItem.ImageUrl,
	}
}

func toDomainRecommendationItem(m mongoRecommedationItem) domain.RecommendationItem {
	return domain.RecommendationItem{
		Text:     m.Text,
		Index:    m.Index,
		Heading:  m.Heading,
		ImageUrl: m.ImageUrl,
	}
}

func toDomainRecommendation(m mongoRecommendation) domain.Recommendation {
	var items []domain.RecommendationItem
	for _, item := range m.Items {
		items = append(items, toDomainRecommendationItem(item))
	}
	return domain.Recommendation{
		ID:         m.ObjectID,
		MetricId:   m.MetricId,
		MetricType: m.MetricType,
		Items:      items,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
