package handler_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A helper function to send request and prepare json response
func makeRequest(t *testing.T, router *mux.Router, method, url string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var requestBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		requestBody = bytes.NewBuffer(b)
	} else {
		requestBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, url, requestBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestIntegration_CreateProduct(t *testing.T) {
	env := setupTestEnv(t)

	// Seeding category first because creation of product requires a valid category
	_, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Clothing"})
	require.NoError(t, err)
	w := makeRequest(t, env.router, http.MethodPost, "/products", domain.Product{
		Name:     "Shirt",
		Price:    900,
		Quantity: 10,
		Brand:    "Zara",
		Category: "Clothing",
	})

	require.Equal(t, http.StatusCreated, w.Code, "body: %s", w.Body.String()) // First checking if status code is the same as expected

	var result domain.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result)) // Encoding the json response to GO struct
	assert.Equal(t, "Shirt", result.Name)                       // Matching names of request Body and response Body
}

func TestIntegration_CreateProduct_MissingName(t *testing.T) {
	env := setupTestEnv(t)

	_, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Clothing"})
	require.NoError(t, err)

	w := makeRequest(t, env.router, http.MethodPost, "/products", domain.Product{
		Price:    900,
		Quantity: 10,
		Brand:    "Zara",
		Category: "Clothing",
	})

	require.Equal(t, http.StatusInternalServerError, w.Code, "expected failure for missing name")

	var handlerErr map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &handlerErr))
	assert.Equal(t, "product name is required", handlerErr["error"])
}

func TestIntegration_CreateProduct_NegativePrice(t *testing.T) {
	env := setupTestEnv(t)

	_, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Clothing"})
	require.NoError(t, err)
	w := makeRequest(t, env.router, http.MethodPost, "/products", domain.Product{
		Name:     "Shirt",
		Price:    -900,
		Quantity: 10,
		Brand:    "Zara",
		Category: "Clothing",
	})

	require.Equal(t, http.StatusInternalServerError, w.Code, "expected failure for negative price")

	var handlerErr map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &handlerErr))
	assert.Equal(t, "product price must be greater than 0", handlerErr["error"])
}

func TestIntegration_CreateProduct_MissingCategory(t *testing.T) {
	env := setupTestEnv(t)

	w := makeRequest(t, env.router, http.MethodPost, "/products", domain.Product{
		Name:     "Shirt",
		Price:    900,
		Quantity: 10,
		Brand:    "Zara",
		Category: "",
	})

	require.Equal(t, http.StatusInternalServerError, w.Code, "expected failure for missing category")

	var handlerErr map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &handlerErr))
	assert.Equal(t, "category name is required", handlerErr["error"])
}

// Given category does not exist in the category collection
func TestIntegration_CreateProduct_InvalidCategory(t *testing.T) {
	env := setupTestEnv(t)

	w := makeRequest(t, env.router, http.MethodPost, "/products", domain.Product{
		Name:     "Shirt",
		Price:    900,
		Quantity: 10,
		Brand:    "Zara",
		Category: "Clothing",
	})

	require.Equal(t, http.StatusInternalServerError, w.Code, "expected failure for invalid category")
	var handlerErr map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &handlerErr))
	assert.Equal(t, "mongo: no documents in result", handlerErr["error"])
}

func TestIntegration_ListProducts(t *testing.T) {
	env := setupTestEnv(t)

	// Seeding category and product collections for testing
	_, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Clothing"})
	require.NoError(t, err)
	_, err = seedCategory(t, env.router, domain.ProductCategory{Title: "Electronics"})
	require.NoError(t, err)

	_, err = seedProduct(t, env.router, domain.Product{
		Name:     "Shirt",
		Price:    900,
		Quantity: 10,
		Brand:    "Zara",
		Category: "Clothing",
	})
	require.NoError(t, err)

	_, err = seedProduct(t, env.router, domain.Product{
		Name:  "Samsung S24",
		Price: 899, Quantity: 5,
		Brand:    "Samsung",
		Category: "Electronics",
	})
	require.NoError(t, err)

	w := makeRequest(t, env.router, http.MethodGet, "/products", nil)
	require.Equal(t, http.StatusOK, w.Code)

	var products []domain.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &products), "failed to decode response")
	assert.Len(t, products, 2) // This confirms that the function is behaving as expected
}

func TestIntegration_GetProduct(t *testing.T) {
	env := setupTestEnv(t)

	_, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Clothing"})
	require.NoError(t, err)
	created, err := seedProduct(t, env.router, domain.Product{
		Name:     "Shirt",
		Price:    900,
		Quantity: 10,
		Brand:    "Zara",
		Category: "Clothing",
	})
	require.NoError(t, err)

	w := makeRequest(t, env.router, http.MethodGet, "/products/"+created.ID, nil)
	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())

	var result domain.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, "Shirt", result.Name)
}

func TestIntegration_FilterProductsByCategory(t *testing.T) {
	env := setupTestEnv(t)

	_, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Electronics"})
	require.NoError(t, err)
	_, err = seedCategory(t, env.router, domain.ProductCategory{Title: "Clothing"})
	require.NoError(t, err)

	_, err = seedProduct(t, env.router, domain.Product{Name: "Shirt", Price: 900, Quantity: 10, Brand: "Zara", Category: "Clothing"})
	require.NoError(t, err)
	_, err = seedProduct(t, env.router, domain.Product{Name: "Samsung S24", Price: 899, Quantity: 5, Brand: "Samsung", Category: "Electronics"})
	require.NoError(t, err)
	_, err = seedProduct(t, env.router, domain.Product{Name: "T-Shirt", Price: 29, Quantity: 100, Brand: "Nike", Category: "Clothing"})
	require.NoError(t, err)

	w := makeRequest(t, env.router, http.MethodGet, "/products/filter/Electronics", nil)
	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())

	var products []domain.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &products))

	assert.Len(t, products, 1)
	for _, p := range products {
		assert.Equal(t, "Electronics", p.Category, "unexpected product in Electronics filter: %s", p.Name) // Checking if all the products in result match the requested category.
	}
}

func TestIntegration_UpdateProduct(t *testing.T) {
	env := setupTestEnv(t)

	_, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Electronics"})
	require.NoError(t, err)

	created, err := seedProduct(t, env.router, domain.Product{Name: "iPhone 15", Price: 999, Quantity: 10, Brand: "Apple", Category: "Electronics"})
	require.NoError(t, err)

	w := makeRequest(t, env.router, http.MethodPut, "/products/"+created.ID, domain.Product{
		Name:     "iPhone 15 Pro",
		Price:    1199,
		Quantity: 8,
		Brand:    "Apple",
		Category: "Electronics",
	})
	assert.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())

	var product domain.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &product))
	assert.Equal(t, product.Price, 1199, "Prices of the result and request do not match.")
	assert.Equal(t, product.Name, "iPhone 15 Pro", "Prices of the result and request do not match.")
}

func TestIntegration_UpdateProduct_InvalidID(t *testing.T) {
	env := setupTestEnv(t)

	w := makeRequest(t, env.router, http.MethodPut, "/products/invalid123", domain.Product{
		Name:     "iPhone 15 Pro",
		Price:    1199,
		Quantity: 8,
		Brand:    "Apple",
		Category: "Electronics",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code, "expected failure for invalid ID")
}

func TestIntegration_DeleteProduct(t *testing.T) {
	env := setupTestEnv(t)

	_, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Electronics"})
	require.NoError(t, err)

	created, err := seedProduct(t, env.router, domain.Product{Name: "iPhone 15", Price: 999, Quantity: 10, Brand: "Apple", Category: "Electronics"})
	require.NoError(t, err)

	w := makeRequest(t, env.router, http.MethodDelete, "/products/"+created.ID, nil)
	require.Equal(t, http.StatusNoContent, w.Code, "body: %s", w.Body.String())

	// Verifying that the product does not exist anymore in collection
	w = makeRequest(t, env.router, http.MethodGet, "/products"+created.ID, nil)
	require.Equal(t, http.StatusNotFound, w.Code, "Product with id: %s still exists in the database", created.ID)
}
