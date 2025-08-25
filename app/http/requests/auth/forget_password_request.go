package auth

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ForgetPasswordRequest struct {
	Mobile      string `form:"mobile" json:"mobile"`
	GapUserName string `form:"gap_user_name" json:"gap_user_name"`
	Code        string `form:"code" json:"code"`
	Password    string `form:"password" json:"password"`
}

func (r *ForgetPasswordRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *ForgetPasswordRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ForgetPasswordRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"gap_user_name": "required",
		"mobile":        "required",
		"code":          "required",
		"password":      "required",
	}
}

func (r *ForgetPasswordRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ForgetPasswordRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ForgetPasswordRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
