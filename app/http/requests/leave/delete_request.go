package leave

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type DeleteRequest struct {
	ID uint `form:"id" json:"id"`
}

func (r *DeleteRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *DeleteRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DeleteRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id": "required",
	}
}

func (r *DeleteRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DeleteRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DeleteRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
