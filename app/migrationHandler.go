package app

import (
	"brokerage/migrations"
	"brokerage/tech"
	"net/http"
)

type MigrationHandler struct {
	service migrations.Migration
}

func (h *MigrationHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := h.service.Migrate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tech.LogWarn("error: cannot create database")
		tech.LogWarn(err.Error())
	}
}

func NewMigrationHandler(service migrations.Migration) *MigrationHandler {
	return &MigrationHandler{
		service,
	}
}
