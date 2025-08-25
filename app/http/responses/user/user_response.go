package user

import "leave/app/models"

type UserResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

func NewUserResource(user models.User) UserResponse {
	return UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Mobile: user.Mobile,
	}
}

type LiteUserResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewLiteUserResource(user models.User) LiteUserResponse {
	return LiteUserResponse{
		ID:   user.ID,
		Name: user.Name,
	}
}

func NewLiteUserCollection(users []models.User) []LiteUserResponse {
	res := make([]LiteUserResponse, len(users))
	for i, u := range users {
		res[i] = NewLiteUserResource(u)
	}
	return res
}
