package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/service"
	"github.com/gorilla/mux"
)

type ProductCategoryHandler struct {
	service *service.ProductCategoryService // a pointer to the ProductCategoryService struct
}

func NewCategoryHandler(categoryService *service.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{service: categoryService} // categoryService pointer passed in main.go after creating ProductCategoryService struct
}

// Creating a new category in collection
func (h *ProductCategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Creating context with timeout of 5 second
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	// Extracting category from request
	var category domain.ProductCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request."})
		return
	}

	// Creating new category by passing it through to the service layer
	resultCategory, err := h.service.CreateCategory(ctx, &category)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong."})
		return
	}

	// setting up response
	sendResponse(w, http.StatusCreated, resultCategory)
}

// Get a slice of all the categories
func (h *ProductCategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	// passing context to service layer
	categories, err := h.service.ListCategories(ctx)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong."})
		return
	}

	sendResponse(w, http.StatusOK, categories)
}

// Updating a particular category with the category ID given in URL
func (h *ProductCategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	// Extracting ID from URL path
	params := mux.Vars(r)
	reqID := params["id"]
	if reqID == "" {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid category ID."})
		return
	}

	// Extracting product from request
	var category domain.ProductCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request."})
		return
	}

	// Passing request to service layer
	if err := h.service.UpdateCategory(ctx, reqID, &category); err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong."})
		return
	}

	sendResponse(w, http.StatusOK, category)
}

// Delete a specific category with it's ID given in URL path
func (h *ProductCategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	// Extracting ID from URL path
	params := mux.Vars(r)
	reqID := params["id"]
	if reqID == "" {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid category ID."})
		return
	}

	// Passing request to service layer
	if err := h.service.DeleteCategory(ctx, reqID); err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong."})
		return
	}
	sendResponse(w, http.StatusNoContent, nil)
}
