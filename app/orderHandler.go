package app

import (
	"brokerage/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	service service.OrderService
}

func (h OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var err error
	orderRes := service.OrderCreateRequest{}

	orderRes.AccountID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("route parameter 'id', was not provided")
		return
	}

	json.NewDecoder(r.Body).Decode(&orderRes.Orders)

	response := h.service.CreateOrder(&orderRes)
	json.NewEncoder(w).Encode(response.CBalance)

}

func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: s,
	}
}
