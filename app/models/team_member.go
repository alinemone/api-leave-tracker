package models

import (
	"github.com/goravel/framework/database/orm"
)

// TeamMember is the pivot linking users to teams (many-to-many) with a manager flag.
type TeamMember struct {
	orm.Model
	TeamID    uint
	UserID    uint
	IsManager bool
	// IsApproved is false while a self-requested membership awaits manager approval.
	IsApproved bool
}
