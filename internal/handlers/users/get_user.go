package handlers

import (
	"errors"
	"net/http"

	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra"
	appErrors "github.com/olad5/AfriHacks2023-stressless-backend/pkg/errors"
	response "github.com/olad5/AfriHacks2023-stressless-backend/pkg/utils"
)

func (u UserHandler) GetLoggedInUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := u.userService.GetLoggedInUser(ctx)
	if err != nil {
		switch {
		case errors.Is(err, infra.ErrUserNotFound):
			response.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		default:
			response.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
			return
		}
	}

	response.SuccessResponse(w, "user retrieved successfully", ToUserDTO(user))
}
