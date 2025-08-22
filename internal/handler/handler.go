package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"order-notification/internal/entity"
	"order-notification/internal/service"
)

// OrderHandler ссылка на бизнес логику
type OrderHandler struct {
	service *service.OrderService
}

// NewOrderHandler конструктор чтобы создать handler
func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// GetOrder ручка получения заказа по UID
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	if !strings.HasPrefix(r.URL.Path, "/order/") {
		respondWithError(w, http.StatusBadRequest, "некорректный путь, нужен /order/<order_uid>")
		return
	}

	orderUIDStr := strings.TrimPrefix(r.URL.Path, "/order/")
	if orderUIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "order_uid обязателен")
		return
	}

	orderUID, err := uuid.Parse(orderUIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "некорректный формат order_uid, нужен UUID")
		return
	}

	order, err := h.service.GetOrderByUID(r.Context(), orderUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

// CreateOrder ручка создания заказа
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC: %v\n", r)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}()

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	var order entity.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		respondWithError(w, http.StatusBadRequest, "не удалось прочитать JSON: "+err.Error())
		return
	}

	err := h.service.SaveOrder(order, r.Context())
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"status":    "ok",
		"order_uid": order.OrderUID.String(),
	})

}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
