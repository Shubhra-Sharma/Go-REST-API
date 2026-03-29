package routes

import (
	"github.com/Shubhra-Sharma/Go-REST-API/internal/handler"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, h *handler.ProductHandler) {
	router.HandleFunc("/products", h.CreateProduct).Methods("POST")
	router.HandleFunc("/products", h.ListProducts).Methods("GET")
	router.HandleFunc("/products/filter/{category}", h.GetProductByCategory).Methods("GET")
	router.HandleFunc("/products/{id}", h.GetProduct).Methods("GET")
	router.HandleFunc("/products/{id}", h.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", h.DeleteProduct).Methods("DELETE")
}
