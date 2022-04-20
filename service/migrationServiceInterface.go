package service

import "brokerage/domain"

type MigrationServiceInterface struct {
	repo *domain.Migration
}

func (m *MigrationServiceInterface) Migrate() {
	m.repo.CreateTables()
}

func NewMigrationServiceInterface(repo *domain.Migration) *MigrationServiceInterface {
	return &MigrationServiceInterface{
		repo,
	}
}
