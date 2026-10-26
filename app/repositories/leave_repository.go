package repositories

import (
	"time"

	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"

	response "leave/app/http/responses/leave"
	"leave/app/interfaces"
	"leave/app/models"
)

type leaveRepository struct{}

func NewLeaveRepository() interfaces.LeaveInterface {
	return &leaveRepository{}
}

func (l leaveRepository) Create(userID uint, leaveType string, startAt time.Time, endAt time.Time, teamID *uint) (models.Leave, error) {
	leave := models.Leave{
		UserId:  userID,
		Type:    leaveType,
		TeamID:  teamID,
		StartAt: startAt,
		EndAt:   endAt,
	}

	if err := facades.Orm().Query().Create(&leave); err != nil {
		return models.Leave{}, err
	}

	return leave, nil
}

// MyLeaves returns the current user's leaves with company/team labels.
func (l leaveRepository) MyLeaves(userID uint, page int, perPage int) ([]response.MyLeaveRow, int64, error) {
	var rows []response.MyLeaveRow
	var total int64

	err := facades.Orm().Query().Table("leaves").
		Select(`
			leaves.id,
			leaves.type,
			COALESCE(teams.name, '') AS team_name,
			leaves.start_at,
			leaves.end_at
		`).
		Join("LEFT JOIN teams ON teams.id = leaves.team_id").
		Where("leaves.user_id = ?", userID).
		Where("leaves.deleted_at IS NULL").
		OrderByDesc("leaves.created_at").
		Paginate(page, perPage, &rows, &total)

	return rows, total, err
}

// MyTotals aggregates the user's own leave days/hours across all their leaves in [from, to).
// Uses Find (not First) so GORM does not append an "ORDER BY id" that would break
// the aggregate (no GROUP BY) query.
func (l leaveRepository) MyTotals(userID uint, from time.Time, to time.Time) (response.AggRow, error) {
	var rows []response.AggRow
	err := facades.Orm().Query().Table("leaves").
		Select(`
			COALESCE(SUM(CASE WHEN leaves.type = 'daily' THEN GREATEST(EXTRACT(EPOCH FROM (leaves.end_at - leaves.start_at)) / 86400, 1) ELSE 0 END), 0) AS total_days,
			COALESCE(SUM(CASE WHEN leaves.type = 'hourly' THEN EXTRACT(EPOCH FROM (leaves.end_at - leaves.start_at)) ELSE 0 END), 0) AS total_hour_seconds
		`).
		Where("leaves.user_id = ?", userID).
		Where("leaves.deleted_at IS NULL").
		Where("leaves.start_at >= ?", from).
		Where("leaves.start_at < ?", to).
		Find(&rows)

	if err != nil || len(rows) == 0 {
		return response.AggRow{}, err
	}
	return rows[0], nil
}

// TeamLeaves returns every leave registered for a team (manager report).
func (l leaveRepository) TeamLeaves(teamID uint, page int, perPage int) ([]response.TeamLeaveRow, int64, error) {
	var rows []response.TeamLeaveRow
	var total int64

	err := facades.Orm().Query().Table("leaves").
		Select(`
			leaves.id,
			users.name AS user_name,
			leaves.type,
			leaves.start_at,
			leaves.end_at
		`).
		Join("JOIN users ON users.id = leaves.user_id").
		Where("leaves.team_id = ?", teamID).
		Where("leaves.deleted_at IS NULL").
		OrderByDesc("leaves.start_at").
		Paginate(page, perPage, &rows, &total)

	return rows, total, err
}

// LeaveTotalsByUser aggregates leave days/hours per user for a team in [from, to).
func (l leaveRepository) LeaveTotalsByUser(teamID uint, from time.Time, to time.Time) ([]response.AggRow, error) {
	var rows []response.AggRow
	err := facades.Orm().Query().Table("leaves").
		Select(`
			leaves.user_id AS user_id,
			COALESCE(SUM(CASE WHEN leaves.type = 'daily' THEN GREATEST(EXTRACT(EPOCH FROM (leaves.end_at - leaves.start_at)) / 86400, 1) ELSE 0 END), 0) AS total_days,
			COALESCE(SUM(CASE WHEN leaves.type = 'hourly' THEN EXTRACT(EPOCH FROM (leaves.end_at - leaves.start_at)) ELSE 0 END), 0) AS total_hour_seconds
		`).
		Where("leaves.team_id = ?", teamID).
		Where("leaves.deleted_at IS NULL").
		Where("leaves.start_at >= ?", from).
		Where("leaves.start_at < ?", to).
		GroupBy("leaves.user_id").
		Find(&rows)
	return rows, err
}

// ResetTeam soft-deletes a team's leaves within [from, to).
func (l leaveRepository) ResetTeam(teamID uint, from time.Time, to time.Time) (int64, error) {
	result, err := facades.Orm().Query().
		Where("team_id = ?", teamID).
		Where("start_at >= ?", from).
		Where("start_at < ?", to).
		Delete(&models.Leave{})
	if err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

// ResetAll soft-deletes every leave within [from, to) (admin-wide reset).
func (l leaveRepository) ResetAll(from time.Time, to time.Time) (int64, error) {
	result, err := facades.Orm().Query().
		Where("start_at >= ?", from).
		Where("start_at < ?", to).
		Delete(&models.Leave{})
	if err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

// FindByUser returns a single leave owned by the user (used before deletion to
// build the cancellation notification).
func (l leaveRepository) FindByUser(userID uint, id uint) (models.Leave, error) {
	var leave models.Leave
	err := facades.Orm().Query().Where("id = ?", id).Where("user_id = ?", userID).First(&leave)
	return leave, err
}

func (l leaveRepository) Delete(userID uint, id uint) error {
	result, err := facades.Orm().Query().
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Delete(&models.Leave{})

	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}
