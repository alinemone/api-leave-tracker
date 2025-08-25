package leave

import (
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
	"strconv"
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

	start, _ := data.Get("start_at")
	end, _ := data.Get("end_at")
	leaveType, _ := data.Get("type")

	startStr, ok := start.(string)
	if !ok {
		return false
	}

	endStr, ok := end.(string)
	if !ok {
		return false
	}

	typeStr, ok := leaveType.(string)
	if !ok {
		return false
	}

	// Check based on leave type
	if typeStr == "daily" {
		// For daily leaves: no other leave should exist on the same day
		count, err := facades.Orm().Query().Table("leaves").
			Where("user_id = ?", userID).
			Where("DATE(start_at) <= DATE(?)", endStr).
			Where("DATE(end_at) >= DATE(?)", startStr).
			Count()

		if err != nil {
			return false
		}

		return count == 0
	} else if typeStr == "hourly" {
		// For hourly leaves: no duplicate leave on same day and overlapping time
		count, err := facades.Orm().Query().Table("leaves").
			Where("user_id = ?", userID).
			Where("start_at < ?", endStr).
			Where("end_at > ?", startStr).
			Count()

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
