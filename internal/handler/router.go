package handler

import (
	"net/http"
)

// NewRouter инициализация нового роутера
func NewRouter(orderHandler *OrderHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/order/", orderHandler.GetOrder)
	mux.HandleFunc("/new_order", orderHandler.CreateOrder)

	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	return mux
}
