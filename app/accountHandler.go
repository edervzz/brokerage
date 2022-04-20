package app

import (
	"brokerage/service"
	"encoding/json"
	"net/http"
)

type AccountHandler struct {
	service service.AccountService
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var request service.AccountCreateRequest

	json.NewDecoder(r.Body).Decode(&request)
	response, appMess := h.service.CreateAccount(&request)
	json.NewEncoder(w).Encode(response)
	if appMess != nil {
		w.WriteHeader(appMess.Code)
	}
}

func NewAccountHandler(s service.AccountService) *AccountHandler {
	return &AccountHandler{
		service: s,
	}
}
