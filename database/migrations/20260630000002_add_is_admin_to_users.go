package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260630000002AddIsAdminToUsers struct{}

func (r *M20260630000002AddIsAdminToUsers) Signature() string {
	return "20260630000002_add_is_admin_to_users"
}

func (r *M20260630000002AddIsAdminToUsers) Up() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.Boolean("is_admin").Default(false)
	})
}

func (r *M20260630000002AddIsAdminToUsers) Down() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropColumn("is_admin")
	})
}
