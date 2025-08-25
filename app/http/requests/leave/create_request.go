package leave

import (
	"fmt"
	"leave/app/helpers"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type CreateRequest struct {
	Type    string    `form:"type" json:"type" example:"daily"`
	StartAt time.Time `form:"start_at" json:"start_at" example:"2025-08-11 14:00:00"`
	EndAt   time.Time `form:"end_at" json:"end_at" example:"2025-08-11 14:00:00"`
}

func (r *CreateRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *CreateRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CreateRequest) Rules(ctx http.Context) map[string]string {
	user, _ := helpers.CurrentUser(ctx)

	return map[string]string{
		"type":     "required",
		"start_at": "required",
		"end_at":   fmt.Sprintf("required|leave_overlap:%d|leave_same_day", user.ID),
	}
}

func (r *CreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CreateRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
