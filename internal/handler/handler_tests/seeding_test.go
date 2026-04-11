package handler_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()
	client, err := database.Connect(context.Background(), "mongodb://localhost:27017", "restapi_test")
	require.NoError(t, err, "failed to connect to test MongoDB")

	// Initializing repositories pointing to test DB
	productRepo := repository.NewMongoProductRepository(client, "restapi_test", "products")
	categoryRepo := repository.NewMongoProductCategoryRepository(client, "restapi_test", "categories")

	// Initializing services
	productService := service.NewProductService(productRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo, productRepo)

	// Initializing handlers
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
func seedCategory(t *testing.T, router *mux.Router, category domain.ProductCategory) (domain.ProductCategory, error) {
	t.Helper() // this declares the corresponding function as a testing function

	body, _ := json.Marshal(category)                                                 // converting Go struct into a JSON byte slice for http request
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body)) // httptest helps to simulate HTTP requests and responses without starting a real server. The mock requests and responses are handled in memory.
	// http.NewRequest expects the request body as an io.Reader, bytes.NewBuffer(body) wraps []byte into a *bytes.Buffer which implements io.Reader.

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		return domain.ProductCategory{}, fmt.Errorf("failed to create category '%s', status: %d, body: %s", category.Title, w.Code, w.Body.String())
	}
	var result domain.ProductCategory
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		return domain.ProductCategory{}, fmt.Errorf("failed to decode created category: %w", err)
	}
	return result, nil
}

// seedProduct creates a product via the API and returns the created product
func seedProduct(t *testing.T, router *mux.Router, product domain.Product) (domain.Product, error) {
	t.Helper()

	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		return domain.Product{}, fmt.Errorf("failed to create product '%s', status: %d, body: %s", product.Name, w.Code, w.Body.String())
	}
	var result domain.Product
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		return domain.Product{}, fmt.Errorf("failed to decode the created product: %w", err)
	}
	return result, nil
}
