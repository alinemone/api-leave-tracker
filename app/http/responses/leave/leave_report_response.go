package user

type LeaveReportResponse struct {
	UserID     uint    `json:"user_id"`
	UserName   string  `json:"user_name"`
	TotalDays  float64 `json:"total_days"`
	TotalHours string  `json:"total_hours"`
}
