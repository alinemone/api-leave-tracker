package repositories

import (
	"leave/app/interfaces"
	"leave/app/models"
	"time"

	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"

	response "leave/app/http/responses/leave"
)

type leaveRepository struct{}

func NewLeaveRepository() interfaces.LeaveInterface {
	return &leaveRepository{}
}

func (l leaveRepository) Report(page int, perPage int) ([]response.LeaveReportResponse, int64, error) {
	var leaves []response.LeaveReportResponse
	var total int64
	err := facades.Orm().Query().Table("leaves").
		Select(`
          leaves.user_id,
          users.name AS user_name,
          SUM(CASE WHEN leaves.type = 'daily' THEN GREATEST(EXTRACT(EPOCH FROM leaves.end_at - leaves.start_at) / 86400, 1) ELSE 0 END) AS total_days,
          CONCAT(
              FLOOR(SUM(CASE WHEN leaves.type = 'hourly' THEN EXTRACT(EPOCH FROM leaves.end_at - leaves.start_at) ELSE 0 END) / 3600),
              ' ساعت ',
              FLOOR((SUM(CASE WHEN leaves.type = 'hourly' THEN EXTRACT(EPOCH FROM leaves.end_at - leaves.start_at) ELSE 0 END) % 3600) / 60),
              ' دقیقه'
          ) AS total_hours
       `).
		Join("JOIN users ON users.id = leaves.user_id").
		GroupBy("leaves.user_id, users.name").
		OrderBy("leaves.user_id").
		Paginate(page, perPage, &leaves, &total)

	return leaves, total, err
}

func (l leaveRepository) GetLeavesByUserID(userID uint, page int, perPage int) ([]models.Leave, int64, error) {
	var leaves []models.Leave
	var total int64
	err := facades.Orm().Query().Table("leaves").
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		OrderByDesc("created_at").
		Paginate(page, perPage, &leaves, &total)

	return leaves, total, err
}

func (l leaveRepository) Create(UserId uint, Type string, StartAt time.Time, EndAt time.Time) (models.Leave, error) {
	leave := models.Leave{
		UserId:  UserId,
		Type:    Type,
		StartAt: StartAt,
		EndAt:   EndAt,
	}

	err := facades.Orm().Query().Create(&leave)
	if err != nil {
		return models.Leave{}, err
	}

	return leave, nil
}

func (l leaveRepository) Delete(UserId uint, ID uint) error {
	result, err := facades.Orm().Query().
		Where("id = ?", ID).
		Where("user_id = ?", UserId).
		Delete(&models.Leave{})

	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}
