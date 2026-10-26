package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260630000006AddCompanyTeamToLeaves struct{}

func (r *M20260630000006AddCompanyTeamToLeaves) Signature() string {
	return "20260630000006_add_company_team_to_leaves"
}

func (r *M20260630000006AddCompanyTeamToLeaves) Up() error {
	return facades.Schema().Table("leaves", func(table schema.Blueprint) {
		table.UnsignedBigInteger("company_id").Nullable()
		table.UnsignedBigInteger("team_id").Nullable()
	})
}

func (r *M20260630000006AddCompanyTeamToLeaves) Down() error {
	return facades.Schema().Table("leaves", func(table schema.Blueprint) {
		table.DropColumn("company_id")
		table.DropColumn("team_id")
	})
}
