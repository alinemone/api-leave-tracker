package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260630000005CreateTeamMembersTable struct{}

func (r *M20260630000005CreateTeamMembersTable) Signature() string {
	return "20260630000005_create_team_members_table"
}

func (r *M20260630000005CreateTeamMembersTable) Up() error {
	if facades.Schema().HasTable("team_members") {
		return nil
	}
	return facades.Schema().Create("team_members", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("team_id")
		table.UnsignedBigInteger("user_id")
		table.Boolean("is_manager").Default(false)
		table.Foreign("team_id").References("id").On("teams").CascadeOnDelete()
		table.Foreign("user_id").References("id").On("users").CascadeOnDelete()
		table.Timestamps()
		table.Unique("team_id", "user_id")
	})
}

func (r *M20260630000005CreateTeamMembersTable) Down() error {
	return facades.Schema().DropIfExists("team_members")
}
