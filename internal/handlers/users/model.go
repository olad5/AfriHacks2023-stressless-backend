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
