package main
import (
	"log"
	"net/http"
	"context"
	"github.com/gorilla/mux"
	"github.com/Shubhra-Sharma/Go-REST-API/routes"
	"github.com/Shubhra-Sharma/Go-REST-API/middlewares"
	"github.com/Shubhra-Sharma/Go-REST-API/database"
    "github.com/joho/godotenv"
)

func main(){
	 err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	//creating context
	ctx := context.Background()

    //calling the database connection function
    db, err := database.Connect(ctx)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    router := mux.NewRouter()
	routes.AllRoutes(router)
    log.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", middlewares.LoggingMiddleware(router)))
}

