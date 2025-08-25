package interfaces

import (
	response "leave/app/http/responses/leave"
	"leave/app/models"
	"time"
)

type LeaveInterface interface {
	Report(page int, perPage int) ([]response.LeaveReportResponse, int64, error)
	GetLeavesByUserID(userID uint, page int, perPage int) ([]models.Leave, int64, error)
	Create(UserId uint, Type string, StartAt time.Time, EndAt time.Time) (models.Leave, error)
	Delete(UserId uint, ID uint) error
}
