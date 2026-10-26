package team

// TeamListRow is a team the current user belongs to, with their role.
type TeamListRow struct {
	TeamID    uint   `json:"team_id"`
	Name      string `json:"name"`
	IsManager bool   `json:"is_manager"`
}

// MemberRow is a member of a team (also used for user search results).
type MemberRow struct {
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	IsManager bool   `json:"is_manager"`
}

// JoinRow is a team the user can request to join (with a pending flag).
type JoinRow struct {
	TeamID  uint   `json:"team_id"`
	Name    string `json:"name"`
	Pending bool   `json:"pending"`
}
