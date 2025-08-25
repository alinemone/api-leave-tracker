package auth

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type LoginRequest struct {
	Mobile   string `form:"mobile" json:"mobile"`
	Password string `form:"password" json:"password"`
}

func (r *LoginRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *LoginRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *LoginRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"mobile":   "required",
		"password": "required",
	}
}

func (r *LoginRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *LoginRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *LoginRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
