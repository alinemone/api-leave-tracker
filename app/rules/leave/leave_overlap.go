package leave

import (
	"leave/app/models"
	"strconv"

	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

const (
	LeaveTypeDaily  = "daily"
	LeaveTypeHourly = "hourly"
)

type LeaveOverlap struct {
}

// Signature The name of the rule.
func (receiver *LeaveOverlap) Signature() string {
	return "leave_overlap"
}

// Passes Determine if the validation rule passes.
func (receiver *LeaveOverlap) Passes(data validation.Data, val any, options ...any) bool {
	if len(options) == 0 {
		return false // userID لازم است
	}

	userIDStr, ok := options[0].(string)
	if !ok {
		return false
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return false
	}

	start, startExists := data.Get("start_at")
	end, endExists := data.Get("end_at")
	if !startExists || !endExists {
		return false
	}

	startStr, ok := start.(string)
	if !ok || startStr == "" {
		return false
	}

	endStr, ok := end.(string)
	if !ok || endStr == "" {
		return false
	}

	leaveType, typeExists := data.Get("type")
	if !typeExists {
		return false
	}

	typeStr, ok := leaveType.(string)
	if !ok {
		return false
	}

	// Check based on leave type
	if typeStr == LeaveTypeDaily {
		// For daily leaves: no other leave should exist on the same day
		query := facades.Orm().Query().Model(&models.Leave{}).
			Where("user_id = ?", userID).
			Where("start_at::date <= (?)::date", endStr).
			Where("end_at::date >= (?)::date", startStr)

		// Exclude current record if it's an update
		if leaveID, exists := data.Get("id"); exists {
			query = query.Where("id != ?", leaveID)
		}

		count, err := query.Count()

		if err != nil {
			return false
		}

		return count == 0
	} else if typeStr == LeaveTypeHourly {
		// For hourly leaves: no duplicate leave on same day and overlapping time
		query := facades.Orm().Query().Model(&models.Leave{}).
			Where("user_id = ?", userID).
			Where("start_at < ?", endStr).
			Where("end_at > ?", startStr)

		// Exclude current record if it's an update
		if leaveID, exists := data.Get("id"); exists {
			query = query.Where("id != ?", leaveID)
		}

		count, err := query.Count()

		if err != nil {
			return false
		}

		return count == 0
	}

	return false
}

// Message Get the validation error message.
func (receiver *LeaveOverlap) Message() string {
	return "شما در این بازه زمانی قبلاً مرخصی ثبت کرده‌اید"
}
