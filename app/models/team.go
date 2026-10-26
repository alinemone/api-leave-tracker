package models

import (
	"github.com/goravel/framework/database/orm"
)

type Team struct {
	orm.Model
	orm.SoftDeletes
	Name string
	// CreatorID is the user who created the team (and is a manager of it).
	CreatorID uint
	// WebhookURL is an optional incoming webhook for team notifications.
	WebhookURL string
}
