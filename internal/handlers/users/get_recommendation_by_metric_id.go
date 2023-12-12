package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra"
	appErrors "github.com/olad5/AfriHacks2023-stressless-backend/pkg/errors"
	response "github.com/olad5/AfriHacks2023-stressless-backend/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func (u UserHandler) GetRecommendationByMetricId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		response.ErrorResponse(w, "id required", http.StatusBadRequest)
		return
	}
	metricType := r.URL.Query().Get("metric_type")
	if metricType == "" {
		response.ErrorResponse(w, "metric_type required", http.StatusBadRequest)
		return
	}

	metricId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrInvalidID.Error(), http.StatusBadRequest)
		return
	}

	recommendation, err := u.userService.GetRecommendationByMetricId(ctx, metricId, metricType)
	if err != nil {
		switch {
		case errors.Is(err, infra.ErrUserNotFound):
			response.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, infra.ErrMetricNotFound):
			response.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, infra.ErrRecommendationNotFound):
			response.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		default:
			u.logger.Error("[internal server error: ]", zap.Error(err))
			response.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
			return
		}
	}

	response.SuccessResponse(w, "recommendation retrieved successfully", ToRecommendationDTO(recommendation))
}
