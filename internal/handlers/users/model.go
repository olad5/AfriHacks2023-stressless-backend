package handlers

import (
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
)

type UserDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func ToUserDTO(user domain.User) UserDTO {
	return UserDTO{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
