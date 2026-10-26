package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260630000003CreateCompaniesTable struct{}

func (r *M20260630000003CreateCompaniesTable) Signature() string {
	return "20260630000003_create_companies_table"
}

func (r *M20260630000003CreateCompaniesTable) Up() error {
	if facades.Schema().HasTable("companies") {
		return nil
	}
	return facades.Schema().Create("companies", func(table schema.Blueprint) {
		table.ID("id")
		table.String("name")
		table.UnsignedBigInteger("owner_id")
		table.Timestamps()
		table.SoftDeletes()
	})
}

func (r *M20260630000003CreateCompaniesTable) Down() error {
	return facades.Schema().DropIfExists("companies")
}
