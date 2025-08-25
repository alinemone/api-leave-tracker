package auth

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ResetCodeRequest struct {
	Mobile      string `form:"mobile" json:"mobile"`
	GapUserName string `form:"gap_user_name" json:"gap_user_name"`
}

func (r *ResetCodeRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *ResetCodeRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ResetCodeRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"gap_user_name": "required",
		"mobile":        "required",
	}
}

func (r *ResetCodeRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ResetCodeRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ResetCodeRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
