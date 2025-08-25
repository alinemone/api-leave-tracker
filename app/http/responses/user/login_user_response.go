package user

import "leave/app/models"

type LoginUserResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func NewLoginResource(user models.User, token string) LoginUserResponse {
	return LoginUserResponse{
		Token: token,
		User:  NewUserResource(user),
	}
}
