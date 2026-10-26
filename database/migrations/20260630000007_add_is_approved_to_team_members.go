package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260630000007AddIsApprovedToTeamMembers struct{}

func (r *M20260630000007AddIsApprovedToTeamMembers) Signature() string {
	return "20260630000007_add_is_approved_to_team_members"
}

func (r *M20260630000007AddIsApprovedToTeamMembers) Up() error {
	return facades.Schema().Table("team_members", func(table schema.Blueprint) {
		// Existing memberships are treated as approved.
		table.Boolean("is_approved").Default(true)
	})
}

func (r *M20260630000007AddIsApprovedToTeamMembers) Down() error {
	return facades.Schema().Table("team_members", func(table schema.Blueprint) {
		table.DropColumn("is_approved")
	})
}
