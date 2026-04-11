package routes

import (
	"github.com/Shubhra-Sharma/Go-REST-API/internal/handler"
	"github.com/gorilla/mux"
)

func CategoryRoutes(router *mux.Router, h handler.ProductCategoryHandlerInterface) {
	router.HandleFunc("/categories", h.ListCategories).Methods("GET")
	router.HandleFunc("/categories", h.CreateCategory).Methods("POST")
	router.HandleFunc("/categories/{id}", h.GetCategoryByID).Methods("GET")
	router.HandleFunc("/categories/{id}", h.UpdateCategory).Methods("PUT")
	router.HandleFunc("/categories/{id}", h.DeleteCategory).Methods("DELETE")
}
