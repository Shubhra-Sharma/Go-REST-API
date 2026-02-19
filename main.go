package main
import (
	"log"
	"net/http"
	"github.com/Shubhra-Sharma/Go-REST-API/model"
	"github.com/gorilla/mux"
	"github.com/Shubhra-Sharma/Go-REST-API/routes"
	"github.com/Shubhra-Sharma/Go-REST-API/middlewares"
)

func main(){
	model.InitializeInventory()
    router := mux.NewRouter()
	routes.AllRoutes(router)
    log.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", middlewares.LoggingMiddleware(router)))
}

