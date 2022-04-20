package app

import (
	"brokerage/service"
	"net/http"
)

type MigrationHandler struct {
	service service.Migration
}

func (h *MigrationHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.service.Migrate()
}

func NewMigrationHandler(service service.Migration) *MigrationHandler {
	return &MigrationHandler{
		service,
	}
}
