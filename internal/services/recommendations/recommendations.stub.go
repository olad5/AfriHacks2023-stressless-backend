package recommendations

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
)

type StubRecommendationService struct {
	client  *http.Client
	baseUrl string
}

type (
	RecommendationService interface {
		GetStresslessScore(ctx context.Context, stressLevel int, mood domain.Mood, sleepQuality domain.SleepQuality, feeling string) (int, error)
		GetRecommendationUsingStressScore(ctx context.Context, metric domain.Metric) (domain.Recommendation, error)
		GetRecommendationUsingStressLevel(ctx context.Context, metric domain.Metric) (domain.Recommendation, error)
		GetRecommendationUsingSleepQuality(ctx context.Context, metric domain.Metric) (domain.Recommendation, error)
		GetRecommendationUsingMood(ctx context.Context, metric domain.Metric) (domain.Recommendation, error)
	}
)

func (s *StubRecommendationService) GetStresslessScore(ctx context.Context, stressLevel int, mood domain.Mood, sleepQuality domain.SleepQuality, feeling string) (int, error) {
	return randomIntWithMaxValueInclusive(20, 95), nil
}

func (s *StubRecommendationService) GetRecommendationUsingStressScore(ctx context.Context, metric domain.Metric) (domain.Recommendation, error) {
	// TODO:TODO: the StressScore recommendation should be 4 and they should be indexed cos you're
	// sending an array back to the user
	// TODO:TODO: sort by index
	return domain.Recommendation{}, nil
}

func (s *StubRecommendationService) GetRecommendationUsingStressLevel(ctx context.Context, metric domain.Metric) (domain.Recommendation, error) {
	return domain.Recommendation{}, nil
}

func (s *StubRecommendationService) GetRecommendationUsingSleepQuality(ctx context.Context, metric domain.Metric) (domain.Recommendation, error) {
	return domain.Recommendation{}, nil
}

func (s *StubRecommendationService) GetRecommendationUsingMood(ctx context.Context, metric domain.Metric) (domain.Recommendation, error) {
	return domain.Recommendation{}, nil
}

// TODO:TODO: this code is duplicated, fix it
func randomIntWithMaxValueInclusive(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
