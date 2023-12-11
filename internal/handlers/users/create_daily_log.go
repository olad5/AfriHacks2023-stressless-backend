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

func (u UserHandler) CreateDailyLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Body == nil {
		response.ErrorResponse(w, appErrors.ErrMissingBody, http.StatusBadRequest)
		return
	}

	type requestDTO struct {
		Mood         domain.Mood         `json:"mood"`
		SleepQuality domain.SleepQuality `json:"sleep_quality"`
		StressLevel  int                 `json:"stress_level"`
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
	newMetric, err := u.userService.CreateDailyLog(ctx, request.StressLevel, domain.Mood(request.Mood), domain.SleepQuality(request.SleepQuality), request.Feeling)
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

	response.SuccessResponse(w, "metric created successfully", ToMetricDTO(newMetric))
}
