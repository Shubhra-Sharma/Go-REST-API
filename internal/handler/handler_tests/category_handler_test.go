package handler_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_CreateCategory(t *testing.T) {
	env := setupTestEnv(t)

	w := makeRequest(t, env.router, http.MethodPost, "/categories", domain.ProductCategory{
		Title:       "Electronics",
		Description: "Electronic devices",
	})

	// Matching status codes
	require.Equal(t, http.StatusCreated, w.Code, "body: %s", w.Body.String())

	var result domain.ProductCategory
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result), "failed to decode response")

	assert.Equal(t, "Electronics", result.Title)
}

func TestIntegration_CreateCategory_MissingTitle(t *testing.T) {
	env := setupTestEnv(t)

	w := makeRequest(t, env.router, http.MethodPost, "/categories", domain.ProductCategory{
		Description: "No title",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code, "expected failure due to missing title")
	var handlerErr map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &handlerErr), "failed to decode error")
	// Checking if the expected error occured or not
	assert.Equal(t, "name of category is required", handlerErr["error"])
}

func TestIntegration_CreateCategory_Duplicate(t *testing.T) {
	env := setupTestEnv(t)

	_, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Electronics"})
	require.NoError(t, err)
	w := makeRequest(t, env.router, http.MethodPost, "/categories", domain.ProductCategory{
		Title: "Electronics",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code, "expected failure for duplicate category")
	var handlerErr map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &handlerErr), "failed to decode error")

	// Checking if the expected error occured or not
	assert.Equal(t, "category 'Electronics' already exists", handlerErr["error"])
}

func TestIntegration_ListCategories(t *testing.T) {
	env := setupTestEnv(t)

	_, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Electronics"})
	require.NoError(t, err)
	_, err = seedCategory(t, env.router, domain.ProductCategory{Title: "Clothing"})
	require.NoError(t, err)
	w := makeRequest(t, env.router, http.MethodGet, "/categories", nil)

	require.Equal(t, http.StatusOK, w.Code)

	var categories []domain.ProductCategory
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &categories), "failed to decode response")

	// Checking if the no of categories is the same as created in seeding or not
	assert.Len(t, categories, 2)
}

func TestIntegration_UpdateCategory(t *testing.T) {
	env := setupTestEnv(t)

	created, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Electronics"})
	require.NoError(t, err)
	w := makeRequest(t, env.router, http.MethodPut, "/categories/"+created.ID, domain.ProductCategory{
		Title:       "Electronics-updated",
		Description: "Updated description",
	})

	assert.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
}

func TestIntegration_UpdateCategory_InvalidID(t *testing.T) {
	env := setupTestEnv(t)

	w := makeRequest(t, env.router, http.MethodPut, "/categories/invalid123", domain.ProductCategory{
		Title: "Updated",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code, "expected failure for invalid ID")
}

func TestIntegration_DeleteCategory(t *testing.T) {
	env := setupTestEnv(t)

	created, err := seedCategory(t, env.router, domain.ProductCategory{Title: "Electronics"})
	require.NoError(t, err)
	w := makeRequest(t, env.router, http.MethodDelete, "/categories/"+created.ID, nil)
	require.Equal(t, http.StatusNoContent, w.Code, "body: %s", w.Body.String())

	// Verifying that the category does not exist anymore in collection
	w = makeRequest(t, env.router, http.MethodGet, "/categories"+created.ID, nil)
	require.Equal(t, http.StatusNotFound, w.Code, "Category with id: %s still exists in the database", created.ID)
}
