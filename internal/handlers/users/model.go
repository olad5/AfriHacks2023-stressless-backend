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
