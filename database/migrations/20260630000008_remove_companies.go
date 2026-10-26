package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260630000008RemoveCompanies struct{}

func (r *M20260630000008RemoveCompanies) Signature() string {
	return "20260630000008_remove_companies"
}

func (r *M20260630000008RemoveCompanies) Up() error {
	if facades.Schema().HasColumn("teams", "company_id") {
		if err := facades.Schema().Table("teams", func(table schema.Blueprint) {
			table.DropColumn("company_id")
		}); err != nil {
			return err
		}
	}
	if facades.Schema().HasColumn("leaves", "company_id") {
		if err := facades.Schema().Table("leaves", func(table schema.Blueprint) {
			table.DropColumn("company_id")
		}); err != nil {
			return err
		}
	}
	return facades.Schema().DropIfExists("companies")
}

func (r *M20260630000008RemoveCompanies) Down() error {
	return facades.Schema().Table("teams", func(table schema.Blueprint) {
		table.UnsignedBigInteger("company_id").Nullable()
	})
}
