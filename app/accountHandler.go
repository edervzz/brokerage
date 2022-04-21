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

	if appMess.Code != 0 {
		w.WriteHeader(appMess.Code)
		json.NewEncoder(w).Encode(appMess.Message)
		return
	}
	json.NewEncoder(w).Encode(response)
}

func NewAccountHandler(s service.AccountService) *AccountHandler {
	return &AccountHandler{
		service: s,
	}
}
