package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Shubhra-Sharma/Go-REST-API/database"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/handler"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/repository"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/service"
	"github.com/Shubhra-Sharma/Go-REST-API/middlewares"
	"github.com/Shubhra-Sharma/Go-REST-API/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DBNAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	catcollectionName := os.Getenv("CATEGORY_COLLECTION_NAME")
	// Checking env variables
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is required")
	}
	if dbName == "" {
		log.Fatal("DBNAME environment variable is required")
	}
	if collectionName == "" {
		collectionName = "products"
	}
	if catcollectionName == "" {
		catcollectionName = "category"
	}
	// creating context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// calling the database connection function
	client, err := database.Connect(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	productRepo := repository.NewMongoProductRepository(client, dbName, collectionName)
	productService := service.NewProductService(productRepo) // MongoDB implementation passed as ProductRepository
	productHandler := handler.NewProductHandler(productService)

	categoryRepo := repository.NewMongoProductCategoryRepository(client, dbName, catcollectionName)
	categoryService := service.NewCategoryService(categoryRepo) // MongoDB implementation passed as ProductRepository
	categoryHandler := handler.NewCategoryHandler(categoryService)

	router := mux.NewRouter()
	routes.RegisterRoutes(router, productHandler)
	routes.CategoryRoutes(router, categoryHandler)
	log.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", middlewares.LoggingMiddleware(router)))
}
