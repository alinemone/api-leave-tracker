package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260630000009AddWebhookUrlToTeams struct{}

func (r *M20260630000009AddWebhookUrlToTeams) Signature() string {
	return "20260630000009_add_webhook_url_to_teams"
}

func (r *M20260630000009AddWebhookUrlToTeams) Up() error {
	return facades.Schema().Table("teams", func(table schema.Blueprint) {
		table.String("webhook_url").Nullable()
	})
}

func (r *M20260630000009AddWebhookUrlToTeams) Down() error {
	return facades.Schema().Table("teams", func(table schema.Blueprint) {
		table.DropColumn("webhook_url")
	})
}
