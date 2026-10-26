package models

import (
	"github.com/goravel/framework/database/orm"
	"time"
)

type Leave struct {
	orm.Model
	orm.SoftDeletes
	UserId uint
	Type   string
	// TeamID is set when the leave is registered for a team (nil = personal).
	TeamID  *uint
	StartAt time.Time
	EndAt   time.Time
}
