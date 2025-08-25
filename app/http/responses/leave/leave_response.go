package user

import (
	"github.com/goravel/framework/support/carbon"
	"gorm.io/gorm"
	"leave/app/models"
	"time"
)

type LeaveResponse struct {
	ID        uint             `json:"id"`
	Type      string           `json:"type"`
	CreatedAt *carbon.DateTime `json:"created_at"`
	DeletedAt gorm.DeletedAt   `json:"deleted_at"`
	StartAt   time.Time        `json:"start_at"`
	EndAt     time.Time        `json:"end_at"`
}

func NewLeaveResource(leave models.Leave) LeaveResponse {
	return LeaveResponse{
		ID:        leave.ID,
		Type:      leave.Type,
		CreatedAt: leave.CreatedAt,
		DeletedAt: leave.DeletedAt,
		StartAt:   leave.StartAt,
		EndAt:     leave.EndAt,
	}
}

func NewLeaveCollection(leaves []models.Leave) []LeaveResponse {
	res := make([]LeaveResponse, len(leaves))
	for i, u := range leaves {
		res[i] = NewLeaveResource(u)
	}
	return res
}
