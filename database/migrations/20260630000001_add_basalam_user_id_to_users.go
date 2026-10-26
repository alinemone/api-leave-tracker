package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260630000001AddBasalamUserIdToUsers struct{}

// Signature The unique signature for the migration.
func (r *M20260630000001AddBasalamUserIdToUsers) Signature() string {
	return "20260630000001_add_basalam_user_id_to_users"
}

// Up Run the migrations.
func (r *M20260630000001AddBasalamUserIdToUsers) Up() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.UnsignedBigInteger("basalam_user_id").Nullable()
		table.Unique("basalam_user_id")
	})
}

// Down Reverse the migrations.
func (r *M20260630000001AddBasalamUserIdToUsers) Down() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropUnique("basalam_user_id")
		table.DropColumn("basalam_user_id")
	})
}
