package models

import (
	"github.com/goravel/framework/database/orm"
	"time"
)

type Leave struct {
	orm.Model
	orm.SoftDeletes
	UserId  uint
	Type    string
	StartAt time.Time
	EndAt   time.Time
}
