package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/services/auth"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/services/recommendations"
)

type UserService struct {
	userRepo              infra.UserRepository
	authService           auth.AuthService
	metricRepo            infra.MetricRepository
	recommendationService recommendations.RecommendationService
	recommendationRepo    infra.RecommendationRepository
	logger                *zap.Logger
}

var (
	ErrUserAlreadyExists    = errors.New("email already exist")
	ErrPasswordIncorrect    = errors.New("invalid credentials")
	ErrInvalidToken         = errors.New("invalid token")
	ErrUserDoesNotOwnMetric = errors.New("user does not own metric")
)

func NewUserService(userRepo infra.UserRepository, authService auth.AuthService, metricRepo infra.MetricRepository, recommendationService recommendations.RecommendationService, recommendationRepo infra.RecommendationRepository, logger *zap.Logger) (*UserService, error) {
	if userRepo == nil {
		return &UserService{}, errors.New("UserService failed to initialize, userRepo is nil")
	}
	if authService == nil {
		return &UserService{}, errors.New("UserService failed to initialize, authService is nil")
	}
	if metricRepo == nil {
		return &UserService{}, errors.New("UserService failed to initialize, metricRepo is nil")
	}
	if recommendationService == nil {
		return &UserService{}, errors.New("UserService failed to initialize, recommendationService is nil")
	}
	if recommendationRepo == nil {
		return &UserService{}, errors.New("UserService failed to initialize, recommendationRepo is nil")
	}
	return &UserService{userRepo, authService, metricRepo, recommendationService, recommendationRepo, logger}, nil
}

func (u *UserService) CreateUser(ctx context.Context, firstName, lastName, email, password string) (domain.User, error) {
	existingUser, err := u.userRepo.GetUserByEmail(ctx, email)
	if err == nil && existingUser.Email == email {
		return domain.User{}, ErrUserAlreadyExists
	}

	hashedPassword, err := hashAndSalt([]byte(password))
	if err != nil {
		return domain.User{}, err
	}

	newUser := domain.User{
		ID:                   primitive.NewObjectID(),
		Email:                email,
		FirstName:            firstName,
		LastName:             lastName,
		Password:             hashedPassword,
		IsOnBoardingComplete: false,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	err = u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return domain.User{}, err
	}
	return newUser, nil
}

func (u *UserService) LogUserIn(ctx context.Context, email, password string) (string, error) {
	existingUser, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if isPasswordCorrect := comparePasswords(existingUser.Password, []byte(password)); !isPasswordCorrect {
		return "", ErrPasswordIncorrect
	}

	accessToken, err := u.authService.GenerateJWT(ctx, existingUser)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (u *UserService) GetLoggedInUser(ctx context.Context) (domain.User, error) {
	jwtClaims, ok := auth.GetJWTClaims(ctx)
	if !ok {
		return domain.User{}, fmt.Errorf("error parsing JWTClaims: %w", ErrInvalidToken)
	}
	userId := jwtClaims.ID

	existingUser, err := u.userRepo.GetUserByUserId(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return existingUser, nil
}

func (u *UserService) CreateDailyLog(ctx context.Context, stressLevel int, mood domain.Mood, sleepQuality domain.SleepQuality, feeling string) (domain.Metric, error) {
	jwtClaims, ok := auth.GetJWTClaims(ctx)
	if !ok {
		return domain.Metric{}, fmt.Errorf("error parsing JWTClaims: %w", ErrInvalidToken)
	}
	userId := jwtClaims.ID

	existingUser, err := u.userRepo.GetUserByUserId(ctx, userId)
	if err != nil {
		return domain.Metric{}, err
	}

	exisitingMetric, err := u.metricRepo.GetUserTodayLogIfExists(ctx, existingUser.ID)
	if err != nil {
		if !errors.Is(err, infra.ErrMetricNotFound) {
			return domain.Metric{}, err
		}
	}
	if exisitingMetric.OwnerId == existingUser.ID {
		return exisitingMetric, nil
	}

	stressLessScore, err := u.recommendationService.GetStresslessScore(ctx, stressLevel, mood, sleepQuality, feeling)
	if err != nil {
		return domain.Metric{}, fmt.Errorf("error generating stressScore: %w", err)
	}

	newMetric := domain.Metric{
		ID:              primitive.NewObjectID(),
		OwnerId:         existingUser.ID,
		StressLevel:     stressLevel,
		Mood:            mood,
		SleepQuality:    sleepQuality,
		StressLessScore: stressLessScore,
		Feeling:         feeling,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = u.metricRepo.CreateMetric(ctx, newMetric)
	if err != nil {
		return domain.Metric{}, err
	}

	err = u.userRepo.UpdateUserLastMetricLog(ctx, existingUser)
	if err != nil {
		return domain.Metric{}, err
	}

	rs, err := generateRecommendations(ctx, u, newMetric)
	if err != nil {
		return domain.Metric{}, err
	}

	for _, recommendation := range rs {
		err = u.recommendationRepo.CreateRecommendation(ctx, recommendation)
		if err != nil {
			return domain.Metric{}, fmt.Errorf("error saving recommendation : %w", err)
		}
	}

	return newMetric, nil
}

func generateRecommendations(ctx context.Context, u *UserService, newMetric domain.Metric) ([]domain.Recommendation, error) {
	stresslessRecommendation, err := u.recommendationService.GetRecommendationUsingStressScore(ctx, newMetric)
	if err != nil {
		return []domain.Recommendation{}, fmt.Errorf("error generating stressless recommendation : %w", err)
	}
	stresslevelRecommendation, err := u.recommendationService.GetRecommendationUsingStressLevel(ctx, newMetric)
	if err != nil {
		return []domain.Recommendation{}, fmt.Errorf("error generating stresslevel recommendation : %w", err)
	}
	stressQualityRecommendation, err := u.recommendationService.GetRecommendationUsingSleepQuality(ctx, newMetric)
	if err != nil {
		return []domain.Recommendation{}, fmt.Errorf("error generating stressQuality recommendation : %w", err)
	}
	moodRecommendation, err := u.recommendationService.GetRecommendationUsingStressLevel(ctx, newMetric)
	if err != nil {
		return []domain.Recommendation{}, fmt.Errorf("error generating mood recommendation : %w", err)
	}
	rs := []domain.Recommendation{
		stresslessRecommendation,
		stresslevelRecommendation,
		stressQualityRecommendation,
		moodRecommendation,
	}
	return rs, nil
}

func (u *UserService) GetMetricByMetricId(ctx context.Context, metricId primitive.ObjectID) (domain.Metric, error) {
	jwtClaims, ok := auth.GetJWTClaims(ctx)
	if !ok {
		return domain.Metric{}, fmt.Errorf("error parsing JWTClaims: %w", ErrInvalidToken)
	}
	userId := jwtClaims.ID

	existingUser, err := u.userRepo.GetUserByUserId(ctx, userId)
	if err != nil {
		return domain.Metric{}, err
	}

	metric, err := u.metricRepo.GetMetricById(ctx, metricId)
	if err != nil {
		return domain.Metric{}, err
	}

	if metric.OwnerId != existingUser.ID {
		return domain.Metric{}, ErrUserDoesNotOwnMetric
	}
	return metric, nil
}

func (u *UserService) GetRecommendationByMetricId(ctx context.Context, metricId primitive.ObjectID, metricType string) (domain.Recommendation, error) {
	jwtClaims, ok := auth.GetJWTClaims(ctx)
	if !ok {
		return domain.Recommendation{}, fmt.Errorf("error parsing JWTClaims: %w", ErrInvalidToken)
	}
	userId := jwtClaims.ID

	existingUser, err := u.userRepo.GetUserByUserId(ctx, userId)
	if err != nil {
		return domain.Recommendation{}, err
	}

	metric, err := u.metricRepo.GetMetricById(ctx, metricId)
	if err != nil {
		return domain.Recommendation{}, err
	}

	if metric.OwnerId != existingUser.ID {
		return domain.Recommendation{}, ErrUserDoesNotOwnMetric
	}
	recommendation, err := u.recommendationRepo.GetRecommendationByMetricId(ctx, metricId, metricType)
	if err != nil {
		return domain.Recommendation{}, err
	}

	return recommendation, nil
}

func (u *UserService) GetMetricForToday(ctx context.Context) (domain.Metric, error) {
	jwtClaims, ok := auth.GetJWTClaims(ctx)
	if !ok {
		return domain.Metric{}, fmt.Errorf("error parsing JWTClaims: %w", ErrInvalidToken)
	}
	userId := jwtClaims.ID

	existingUser, err := u.userRepo.GetUserByUserId(ctx, userId)
	if err != nil {
		return domain.Metric{}, err
	}

	metric, err := u.metricRepo.GetUserTodayLogIfExists(ctx, userId)
	if err != nil {
		return domain.Metric{}, err
	}

	if metric.OwnerId != existingUser.ID {
		return domain.Metric{}, ErrUserDoesNotOwnMetric
	}
	return metric, nil
}

func (u *UserService) GetRecentMetricsByUserId(ctx context.Context) ([]domain.Metric, error) {
	// TODO:TODO: this method has issues
	jwtClaims, ok := auth.GetJWTClaims(ctx)
	if !ok {
		return []domain.Metric{}, fmt.Errorf("error parsing JWTClaims: %w", ErrInvalidToken)
	}
	userId := jwtClaims.ID

	// TODO:TODO: I should add the limit, rowperpage and offset stuff
	existingUser, err := u.userRepo.GetUserByUserId(ctx, userId)
	if err != nil {
		return []domain.Metric{}, err
	}
	metrics, err := u.metricRepo.GetRecentMetricsByUserId(ctx, existingUser.ID)
	if err != nil {
		return []domain.Metric{}, err
	}
	return metrics, nil
}

func (u *UserService) CompleteUserOnboarding(ctx context.Context, stressLevel int, mood domain.Mood, sleepQuality domain.SleepQuality, feeling string) (domain.User, error) {
	jwtClaims, ok := auth.GetJWTClaims(ctx)
	if !ok {
		return domain.User{}, fmt.Errorf("error parsing JWTClaims: %w", ErrInvalidToken)
	}
	userId := jwtClaims.ID

	existingUser, err := u.userRepo.GetUserByUserId(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	if existingUser.IsOnBoardingComplete {
		return existingUser, err
	}

	stressLessScore, err := u.recommendationService.GetStresslessScore(ctx, stressLevel, mood, sleepQuality, feeling)
	if err != nil {
		return domain.User{}, fmt.Errorf("error generating stressScore: %w", err)
	}

	newMetric := domain.Metric{
		ID:              primitive.NewObjectID(),
		OwnerId:         userId,
		StressLevel:     stressLevel,
		Mood:            mood,
		SleepQuality:    sleepQuality,
		StressLessScore: stressLessScore,
		Feeling:         feeling,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = u.metricRepo.CreateMetric(ctx, newMetric)
	if err != nil {
		return domain.User{}, err
	}

	// TODO:TODO: this url might help  https://frontendmasters.com/courses/openai-node/
	rs, err := generateRecommendations(ctx, u, newMetric)
	if err != nil {
		return domain.User{}, err
	}

	for _, recommendation := range rs {
		err = u.recommendationRepo.CreateRecommendation(ctx, recommendation)
		if err != nil {
			return domain.User{}, fmt.Errorf("error saving recommendation : %w", err)
		}
	}

	updatedUser := domain.User{
		ID:                   existingUser.ID,
		Email:                existingUser.Email,
		FirstName:            existingUser.FirstName,
		LastName:             existingUser.LastName,
		Password:             existingUser.Password,
		IsOnBoardingComplete: true,
		LastMetricLog:        time.Now(),
		UpdatedAt:            time.Now(),
	}
	err = u.userRepo.UpdateUser(ctx, updatedUser)
	if err != nil {
		return domain.User{}, err
	}
	return updatedUser, nil
}

func hashAndSalt(plainPassword []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(plainPassword, bcrypt.MinCost)
	if err != nil {
		return "", errors.New("error hashing password")
	}
	return string(hash), nil
}

func comparePasswords(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	return err == nil
}
