package migrations

import "brokerage/domain"

type MigrationServiceInterface struct {
	repo *domain.Migration
}

func (m *MigrationServiceInterface) Migrate() error {
	return m.repo.CreateTables()
}

func NewMigrationServiceInterface(repo *domain.Migration) *MigrationServiceInterface {
	return &MigrationServiceInterface{
		repo,
	}
}
