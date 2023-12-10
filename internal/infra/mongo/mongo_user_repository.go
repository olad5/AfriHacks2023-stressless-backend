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

var contextTimeoutDuration = 5 * time.Second

type MongoUserRepository struct {
	users  *mongo.Collection
	logger *zap.Logger
}

func NewMongoUserRepo(ctx context.Context, mongoDatabase *mongo.Database, logger *zap.Logger) (*MongoUserRepository, error) {
	userCollection := mongoDatabase.Collection("users")

	return &MongoUserRepository{users: userCollection, logger: logger}, nil
}

func (m *MongoUserRepository) CreateUser(ctx context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoUser := toMongoUser(user)

	_, err := m.users.InsertOne(ctx, mongoUser)
	if err != nil {
		m.logger.Error("failed to persist user: %w", zap.Error(err))
		return fmt.Errorf("failed to persist user: %w", err)
	}
	return nil
}

func (m *MongoUserRepository) UpdateUser(ctx context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoUser := toMongoUser(user)

	filter := bson.M{"_id": user.ID}
	updatedDoc := bson.M{
		"$set": mongoUser,
	}
	_, err := m.users.UpdateOne(ctx, filter, updatedDoc)
	if err != nil {
		m.logger.Error("failed to update user: %w", zap.Error(err))
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (m *MongoUserRepository) GetUserByEmail(ctx context.Context, userEmail string) (domain.User, error) {
	user := mongoUser{}
	err := m.users.FindOne(ctx, bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		m.logger.Error("failed retrieve user by email: %w", zap.Error(err))
		return domain.User{}, infra.ErrUserNotFound
	}
	return toDomainUser(user), nil
}

func (m *MongoUserRepository) GetUserByUserId(ctx context.Context, userId primitive.ObjectID) (domain.User, error) {
	user := mongoUser{}
	err := m.users.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		m.logger.Error("failed retrieve user by id: %w", zap.Error(err))
		return domain.User{}, infra.ErrUserNotFound
	}

	return toDomainUser(user), nil
}

type mongoUser struct {
	ObjectID            primitive.ObjectID `bson:"_id"`
	Email               string             `bson:"email"`
	FirstName           string             `bson:"first_name"`
	LastName            string             `bson:"last_name"`
	Password            string             `bson:"password"`
	IsOnBoardinComplete bool               `bson:"is_onboarding_complete"`
	LastMetricLog       time.Time          `bson:"last_metric_log"`
	CreatedAt           time.Time          `bson:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at"`
}

func toMongoUser(user domain.User) mongoUser {
	return mongoUser{
		ObjectID:            user.ID,
		Email:               user.Email,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Password:            user.Password,
		IsOnBoardinComplete: user.IsOnBoardingComplete,
		LastMetricLog:       user.LastMetricLog,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}

func toDomainUser(m mongoUser) domain.User {
	return domain.User{
		ID:                   m.ObjectID,
		Email:                m.Email,
		FirstName:            m.FirstName,
		LastName:             m.LastName,
		Password:             m.Password,
		LastMetricLog:        m.LastMetricLog,
		IsOnBoardingComplete: m.IsOnBoardinComplete,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}
