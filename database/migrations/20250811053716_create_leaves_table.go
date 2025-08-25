package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250811053716CreateLeavesTable struct{}

// Signature The unique signature for the migration.
func (r *M20250811053716CreateLeavesTable) Signature() string {
	return "20250811053716_create_leaves_table"
}

// Up Run the migrations.
func (r *M20250811053716CreateLeavesTable) Up() error {
	if !facades.Schema().HasTable("leaves") {
		return facades.Schema().Create("leaves", func(table schema.Blueprint) {
			table.ID("id")
			table.Integer("user_id")
			table.Foreign("user_id").References("id").On("users").CascadeOnDelete()
			table.Enum("type", []any{"daily", "hourly"})
			table.DateTimeTz("start_at")
			table.DateTimeTz("end_at")
			table.TimestampsTz()
			table.SoftDeletes()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250811053716CreateLeavesTable) Down() error {
	return facades.Schema().DropIfExists("leaves")
}
