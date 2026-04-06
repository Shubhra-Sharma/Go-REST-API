package handler_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Shubhra-Sharma/Go-REST-API/database"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/handler"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/repository"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/service"
	"github.com/Shubhra-Sharma/Go-REST-API/routes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type testEnv struct {
	router          *mux.Router
	productHandler  handler.ProductHandlerInterface
	categoryHandler handler.ProductCategoryHandlerInterface
	productRepo     repository.ProductRepository
	categoryRepo    repository.ProductCategoryRepository
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// // Loading testDB url
	// uri := os.Getenv("MONGO_TEST_URI")
	// if uri == "" {
	// 	log.Fatal("MONGO_TEST_URI environment variable is required.")
	// }
	// Connecting to local MongoDB running in Docker
	client, err := database.Connect(context.Background(), "mongodb://localhost:27017", "restapi_test")
	require.NoError(t, err, "failed to connect to test MongoDB")

	// Initialize repositories pointing to test DB
	productRepo := repository.NewMongoProductRepository(client, "restapi_test", "products")
	categoryRepo := repository.NewMongoProductCategoryRepository(client, "restapi_test", "categories")

	// Initialize services
	productService := service.NewProductService(productRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo, productRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Setup router with all routes
	router := mux.NewRouter()
	routes.RegisterRoutes(router, productHandler)
	routes.CategoryRoutes(router, categoryHandler)

	setupEnv := &testEnv{
		router:          router,
		productHandler:  productHandler,
		categoryHandler: categoryHandler,
		productRepo:     productRepo,
		categoryRepo:    categoryRepo,
	}

	// cleaning test DB before each test
	cleanDB(t, client)

	// cleaning test DB after each test
	t.Cleanup(func() { // Cleanup takes a function which will be called when the test and all its subtests complete.
		cleanDB(t, client)
		client.Disconnect(context.Background())
	})

	return setupEnv
}

// cleanDB wipes both collections so each test starts fresh
func cleanDB(t *testing.T, client *mongo.Client) {
	t.Helper()
	ctx := context.Background()

	_, err := client.Database("restapi_test").Collection("products").DeleteMany(ctx, bson.M{})
	require.NoError(t, err, "Could not clean products collection inside cleanDB")

	_, err = client.Database("restapi_test").Collection("categories").DeleteMany(ctx, bson.M{})
	require.NoError(t, err, "Could not clean category collection inside cleanDB")
}

// seedCategory creates a category via the API and returns the created category
func seedCategory(t *testing.T, router *mux.Router, category domain.ProductCategory) domain.ProductCategory {
	t.Helper() // this declares the corresponding function as a testing function

	body, _ := json.Marshal(category)                                                 // converting Go struct into a JSON byte slice for http request
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body)) // httptest helps to simulate HTTP requests and responses without starting a real server. The mock requests and responses are handled in memory.
	// http.NewRequest expects the request body as an io.Reader, bytes.NewBuffer(body) wraps []byte into a *bytes.Buffer which implements io.Reader.

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code,
		"seed: failed to create category '%s', body: %s", category.Title, w.Body.String())

	var result domain.ProductCategory
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result), // converting the JSON string back to GO struct
		"seed: failed to decode created category")

	return result
}

// seedProduct creates a product via the API and returns the created product
func seedProduct(t *testing.T, router *mux.Router, product domain.Product) domain.Product {
	t.Helper()

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code,
		"seed: failed to create product '%s', body: %s", product.Name, w.Body.String())

	var result domain.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result),
		"seed: failed to decode created product")

	return result
}
