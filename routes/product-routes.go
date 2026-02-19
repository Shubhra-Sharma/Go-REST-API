package routes

import(
	"github.com/gorilla/mux"
	"github.com/Shubhra-Sharma/Go-REST-API/handlers"
)

func AllRoutes(router *mux.Router){
	router.HandleFunc("/products",handlers.ProductsHandler).Methods("GET")
	router.HandleFunc("/products/{id}",handlers.ProductHandler).Methods("GET")
	router.HandleFunc("/products",handlers.AddProductHandler).Methods("POST")
	router.HandleFunc("/products/{id}",handlers.UpdateProductHandler).Methods("PUT")
	router.HandleFunc("/products/{id}",handlers.DeleteProductHandler).Methods("DELETE")
}