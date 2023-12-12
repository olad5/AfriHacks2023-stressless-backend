package handlers

import (
	"errors"
	"net/http"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/usecases/users"
	appErrors "github.com/olad5/AfriHacks2023-stressless-backend/pkg/errors"
	response "github.com/olad5/AfriHacks2023-stressless-backend/pkg/utils"
	"go.uber.org/zap"
)

func (u UserHandler) GetMetricForToday(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metric, err := u.userService.GetMetricForToday(ctx)
	if err != nil {
		switch {
		case errors.Is(err, infra.ErrUserNotFound):
			response.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, users.ErrUserDoesNotOwnMetric):
			response.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		default:
			u.logger.Error("[internal server error: ]", zap.Error(err))
			response.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
			return
		}
	}

	response.SuccessResponse(w, "metric retrieved successfully", ToMetricDTO(metric))
}
