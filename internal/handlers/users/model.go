package handlers

import (
	"time"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
)

type UserDTO struct {
	ID                   string     `json:"id"`
	Email                string     `json:"email"`
	FirstName            string     `json:"first_name"`
	LastName             string     `json:"last_name"`
	IsOnBoardingComplete bool       `json:"is_onboarding_complete"`
	LastMetricLog        *time.Time `json:"last_metric_log,omitempty"`
}

func ToUserDTO(user domain.User) UserDTO {
	if user.LastMetricLog.IsZero() {
		return UserDTO{
			ID:                   user.ID.Hex(),
			Email:                user.Email,
			FirstName:            user.FirstName,
			LastName:             user.LastName,
			IsOnBoardingComplete: user.IsOnBoardingComplete,
		}
	}
	return UserDTO{
		ID:                   user.ID.Hex(),
		Email:                user.Email,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		IsOnBoardingComplete: user.IsOnBoardingComplete,
		LastMetricLog:        &user.LastMetricLog,
	}
}

type MetricDTO struct {
	ID              string     `json:"id"`
	OwnerId         string     `json:"owner_id"`
	StressLevel     int        `json:"stress_level"`
	Mood            string     `json:"mood"`
	SleepQuality    string     `json:"sleep_quality"`
	StressLessScore int        `json:"stress_less_score"`
	Feeling         string     `json:"feeling"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

func ToMetricDTO(metric domain.Metric) MetricDTO {
	return MetricDTO{
		ID:              metric.ID.Hex(),
		OwnerId:         metric.OwnerId.Hex(),
		StressLevel:     metric.StressLevel,
		Mood:            string(metric.Mood),
		SleepQuality:    string(metric.SleepQuality),
		StressLessScore: metric.StressLessScore,
		Feeling:         metric.Feeling,
		CreatedAt:       &metric.CreatedAt,
		UpdatedAt:       &metric.UpdatedAt,
	}
}

// ----------------------------------
// stress less scores start
// ----------------------------------
type StatsStressLessScoreDTO struct {
	MetricId        string     `json:"metric_id"`
	StressLessScore int        `json:"stress_less_score"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
type StatsStressLessScorePagedDTO struct {
	// TODO:TODO: i am not sure about this limit, i think it should be page
	Limit int                       `json:"limit"`
	Items []StatsStressLessScoreDTO `json:"items"`
}

func ToStatsStressLessScoreDTO(metric domain.Metric) StatsStressLessScoreDTO {
	return StatsStressLessScoreDTO{
		MetricId:        metric.ID.Hex(),
		StressLessScore: metric.StressLessScore,
		CreatedAt:       &metric.CreatedAt,
		UpdatedAt:       &metric.UpdatedAt,
	}
}

func ToStatsStressLessScorePagedDTO(metrics []domain.Metric) StatsStressLessScorePagedDTO {
	items := []StatsStressLessScoreDTO{}
	for _, metric := range metrics {
		items = append(items, ToStatsStressLessScoreDTO(metric))
	}
	return StatsStressLessScorePagedDTO{
		Limit: len(items),
		Items: items,
	}
}

// ----------------------------------
// mood stats start
// ----------------------------------
type StatsMoodDTO struct {
	MetricId  string     `json:"metric_id"`
	Mood      string     `json:"mood"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
type StatsMoodPagedDTO struct {
	// TODO:TODO: i am not sure about this limit, i think it should be page
	Limit int            `json:"limit"`
	Items []StatsMoodDTO `json:"items"`
}

func ToStatsMoodDTO(metric domain.Metric) StatsMoodDTO {
	return StatsMoodDTO{
		MetricId:  metric.ID.Hex(),
		Mood:      string(metric.Mood),
		CreatedAt: &metric.CreatedAt,
		UpdatedAt: &metric.UpdatedAt,
	}
}

func ToStatsMoodPagedDTO(metrics []domain.Metric) StatsMoodPagedDTO {
	items := []StatsMoodDTO{}
	for _, metric := range metrics {
		items = append(items, ToStatsMoodDTO(metric))
	}
	return StatsMoodPagedDTO{
		Limit: len(items),
		Items: items,
	}
}

// ----------------------------------
// sleep_quality  stats start
// ----------------------------------
type StatsSleepQualityDTO struct {
	MetricId     string     `json:"metric_id"`
	SleepQuality string     `json:"sleep_quality"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
type StatsSleepQualityPagedDTO struct {
	// TODO:TODO: i am not sure about this limit, i think it should be page
	Limit int                    `json:"limit"`
	Items []StatsSleepQualityDTO `json:"items"`
}

func ToStatsSleepQualityDTO(metric domain.Metric) StatsSleepQualityDTO {
	return StatsSleepQualityDTO{
		MetricId:     metric.ID.Hex(),
		SleepQuality: string(metric.SleepQuality),
		CreatedAt:    &metric.CreatedAt,
		UpdatedAt:    &metric.UpdatedAt,
	}
}

func ToStatsSleepQualityPagedDTO(metrics []domain.Metric) StatsSleepQualityPagedDTO {
	items := []StatsSleepQualityDTO{}
	for _, metric := range metrics {
		items = append(items, ToStatsSleepQualityDTO(metric))
	}
	return StatsSleepQualityPagedDTO{
		Limit: len(items),
		Items: items,
	}
}

type RecommendationItemDTO struct {
	Index    int    `json:"index"`
	Heading  string `json:"heading"`
	Text     string `json:"text"`
	ImageUrl string `json:"image_url"`
}

type RecommendationDTO struct {
	ID         string                  `json:"id"`
	MetricId   string                  `json:"metric_id"`
	MetricType string                  `json:"metric_type"`
	Items      []RecommendationItemDTO `json:"items"`
	CreatedAt  *time.Time              `json:"created_at"`
	UpdatedAt  *time.Time              `json:"updated_at"`
}

func ToRecommendationItemDTO(r domain.RecommendationItem) RecommendationItemDTO {
	return RecommendationItemDTO{
		Index:    r.Index,
		Heading:  r.Heading,
		Text:     r.Text,
		ImageUrl: r.ImageUrl,
	}
}

func ToRecommendationDTO(recommendation domain.Recommendation) RecommendationDTO {
	items := []RecommendationItemDTO{}
	for _, metric := range recommendation.Items {
		items = append(items, ToRecommendationItemDTO(metric))
	}
	return RecommendationDTO{
		ID:         recommendation.ID.Hex(),
		MetricId:   recommendation.MetricId.Hex(),
		MetricType: recommendation.MetricType,
		Items:      items,
		CreatedAt:  &recommendation.CreatedAt,
		UpdatedAt:  &recommendation.UpdatedAt,
	}
}
