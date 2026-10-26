package interfaces

import (
	"time"

	response "leave/app/http/responses/leave"
	"leave/app/models"
)

type LeaveInterface interface {
	Create(userID uint, leaveType string, startAt time.Time, endAt time.Time, teamID *uint) (models.Leave, error)
	MyLeaves(userID uint, page int, perPage int) ([]response.MyLeaveRow, int64, error)
	MyTotals(userID uint, from time.Time, to time.Time) (response.AggRow, error)
	TeamLeaves(teamID uint, page int, perPage int) ([]response.TeamLeaveRow, int64, error)
	LeaveTotalsByUser(teamID uint, from time.Time, to time.Time) ([]response.AggRow, error)
	ResetTeam(teamID uint, from time.Time, to time.Time) (int64, error)
	ResetAll(from time.Time, to time.Time) (int64, error)
	FindByUser(userID uint, id uint) (models.Leave, error)
	Delete(userID uint, id uint) error
}
