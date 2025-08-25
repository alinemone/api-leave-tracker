package auth

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type RegisterRequest struct {
	Name     string `form:"name" json:"name"`
	IsActive bool   `form:"is_active" json:"is_active"`
	Mobile   string `form:"mobile" json:"mobile"`
	Password string `form:"password" json:"password"`
}

func (r *RegisterRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *RegisterRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *RegisterRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":      "required",
		"is_active": "bool|required",
		"mobile":    "required",
		"password":  "required",
	}
}

func (r *RegisterRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *RegisterRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *RegisterRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
