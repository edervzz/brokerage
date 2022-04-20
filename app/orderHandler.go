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
	orderRes := service.OrderCreateRequest{}
	orders := []service.OrderCreate{}

	accountId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		json.NewEncoder(w).Encode("cannot get account id")
		w.WriteHeader(http.StatusBadRequest)
	}

	orderRes.AccountID = accountId
	json.NewDecoder(r.Body).Decode(&orders)
	orderRes.Orders = orders

	response := h.service.CreateOrder(&orderRes)
	json.NewEncoder(w).Encode(response)

}

func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: s,
	}
}
