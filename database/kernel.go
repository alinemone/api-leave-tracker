package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"

	"leave/database/migrations"
	"leave/database/seeders"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateUsersTable{},
		&migrations.M20210101000002CreateJobsTable{},
		&migrations.M20250811053716CreateLeavesTable{},
		&migrations.M20260630000001AddBasalamUserIdToUsers{},
		&migrations.M20260630000002AddIsAdminToUsers{},
		&migrations.M20260630000003CreateCompaniesTable{},
		&migrations.M20260630000004CreateTeamsTable{},
		&migrations.M20260630000005CreateTeamMembersTable{},
		&migrations.M20260630000006AddCompanyTeamToLeaves{},
		&migrations.M20260630000007AddIsApprovedToTeamMembers{},
		&migrations.M20260630000008RemoveCompanies{},
		&migrations.M20260630000009AddWebhookUrlToTeams{},
	}
}

func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
	}
}
