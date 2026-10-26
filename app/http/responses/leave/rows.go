package leave

import "time"

// MyLeaveRow is a row in the current user's leaves list (with target label).
type MyLeaveRow struct {
	ID       uint      `json:"id"`
	Type     string    `json:"type"`
	TeamName string    `json:"team_name"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}

// TeamLeaveRow is a row in a team's leaves report (manager view).
type TeamLeaveRow struct {
	ID       uint      `json:"id"`
	UserName string    `json:"user_name"`
	Type     string    `json:"type"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}

// AggRow is a per-user leave aggregate over a period (used by the dashboard).
type AggRow struct {
	UserID           uint    `json:"user_id"`
	TotalDays        float64 `json:"total_days"`
	TotalHourSeconds float64 `json:"total_hour_seconds"`
}
