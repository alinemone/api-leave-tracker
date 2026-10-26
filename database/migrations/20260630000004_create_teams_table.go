package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260630000004CreateTeamsTable struct{}

func (r *M20260630000004CreateTeamsTable) Signature() string {
	return "20260630000004_create_teams_table"
}

func (r *M20260630000004CreateTeamsTable) Up() error {
	if facades.Schema().HasTable("teams") {
		return nil
	}
	return facades.Schema().Create("teams", func(table schema.Blueprint) {
		table.ID("id")
		table.String("name")
		table.UnsignedBigInteger("company_id")
		table.UnsignedBigInteger("creator_id")
		table.Foreign("company_id").References("id").On("companies").CascadeOnDelete()
		table.Timestamps()
		table.SoftDeletes()
	})
}

func (r *M20260630000004CreateTeamsTable) Down() error {
	return facades.Schema().DropIfExists("teams")
}
