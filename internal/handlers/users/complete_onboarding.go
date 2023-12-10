package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra"
	appErrors "github.com/olad5/AfriHacks2023-stressless-backend/pkg/errors"
	response "github.com/olad5/AfriHacks2023-stressless-backend/pkg/utils"
	"go.uber.org/zap"
)

func (u UserHandler) CompleteOnboarding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Body == nil {
		response.ErrorResponse(w, appErrors.ErrMissingBody, http.StatusBadRequest)
		return
	}

	type requestDTO struct {
		Mood         domain.Mood         `json:"mood"`
		SleepQuality domain.SleepQuality `json:"sleep_quality"`
		StressLevel  int64               `json:"stress_level"`
		Feeling      string              `json:"feeling"`
	}
	var request requestDTO
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrInvalidJson, http.StatusBadRequest)
		return
	}
	if !isValidMood(request.Mood) {
		response.ErrorResponse(w, "invalid mood", http.StatusBadRequest)
		return
	}

	if !isValidSleepQuality(request.SleepQuality) {
		response.ErrorResponse(w, "invalid sleep_quality", http.StatusBadRequest)
		return
	}

	if request.StressLevel == 0 {
		response.ErrorResponse(w, "stress_level required", http.StatusBadRequest)
		return
	}
	if request.StressLevel < 0 {
		response.ErrorResponse(w, "stress_level must be a non-negative integer", http.StatusBadRequest)
		return
	}
	user, err := u.userService.CompleteUserOnboarding(ctx, request.StressLevel, domain.Mood(request.Mood), domain.SleepQuality(request.SleepQuality), request.Feeling)
	if err != nil {
		switch {
		case errors.Is(err, infra.ErrUserNotFound):
			response.ErrorResponse(w, "user does not exist", http.StatusNotFound)
			return
		default:
			u.logger.Error("[internal server error: ]", zap.Error(err))
			response.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
			return
		}
	}

	response.SuccessResponse(w, "user onboarding completed successfully", ToUserDTO(user))
}

func isValidMood(m domain.Mood) bool {
	validMoods := map[domain.Mood]bool{
		domain.OVERJOYED: true,
		domain.HAPPY:     true,
		domain.NEUTRAL:   true,
		domain.SAD:       true,
		domain.DEPRESSED: true,
	}
	return validMoods[m]
}

func isValidSleepQuality(q domain.SleepQuality) bool {
	validSleepQualities := map[domain.SleepQuality]bool{
		domain.EXCELLENT: true,
		domain.GOOD:      true,
		domain.FAIR:      true,
		domain.POOR:      true,
		domain.WORST:     true,
	}
	return validSleepQualities[q]
}
