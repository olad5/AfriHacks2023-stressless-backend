package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/usecases/users"
	appErrors "github.com/olad5/AfriHacks2023-stressless-backend/pkg/errors"
	response "github.com/olad5/AfriHacks2023-stressless-backend/pkg/utils"
	"go.uber.org/zap"
)

func (u UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Body == nil {
		response.ErrorResponse(w, appErrors.ErrMissingBody, http.StatusBadRequest)
		return
	}
	type requestDTO struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}
	var request requestDTO
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrInvalidJson, http.StatusBadRequest)
		return
	}
	if request.Email == "" {
		response.ErrorResponse(w, "email required", http.StatusBadRequest)
		return
	}
	if request.Password == "" {
		response.ErrorResponse(w, "password required", http.StatusBadRequest)
		return
	}

	if request.FirstName == "" {
		response.ErrorResponse(w, "first_name required", http.StatusBadRequest)
		return
	}
	if request.LastName == "" {
		response.ErrorResponse(w, "last_name required", http.StatusBadRequest)
		return
	}

	newUser, err := u.userService.CreateUser(ctx, request.FirstName, request.LastName, request.Email, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserAlreadyExists):
			response.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		default:
			u.logger.Error("[internal server error: ]", zap.Error(err))
			response.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
			return
		}
	}

	response.SuccessResponse(w, "user created successfully", ToUserDTO(newUser))
}
