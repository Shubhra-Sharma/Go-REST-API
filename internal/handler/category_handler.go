package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/service"
	"github.com/gorilla/mux"
)

type ProductCategoryHandler struct {
	service *service.ProductCategoryService // a pointer to the ProductCategoryService struct
}

func NewCategoryHandler(cat_Service *service.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{service: cat_Service} // cat_Service pointer passed in main.go after creating ProductCategoryService struct
}

// Creating a new Product in collection
func (h *ProductCategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Creating context with timeout of 5 second
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Extracting category from request
	var category domain.ProductCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request."})
		return
	}

	// Creating new category by passing it through to the service layer
	if err := h.service.CreateCategory(ctx, &category); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong."})
		return
	}

	// setting up response
	result, err := json.Marshal(category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to encode response"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

// Get a slice of all the categories
func (h *ProductCategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// passing context to service layer
	categories, err := h.service.ListCategories(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong."})
		return
	}

	result, err := json.Marshal(categories)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to encode response"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Updating a particular category with the category ID given in URL
func (h *ProductCategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Extracting ID from URL path
	params := mux.Vars(r)
	reqID := params["id"]
	if reqID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid category ID"})
		return
	}

	// Extracting product from request
	var category domain.ProductCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request."})
		return
	}

	// Passing request to service layer
	if err := h.service.UpdateCategory(ctx, reqID, &category); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong."})
		return
	}

	result, err := json.Marshal(category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to encode response"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Delete a specific category with it's ID given in URL path
func (h *ProductCategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Extracting ID from URL path
	params := mux.Vars(r)
	reqID := params["id"]
	if reqID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid product ID"})
		return
	}

	// Passing request to service layer
	if err := h.service.DeleteCategory(ctx, reqID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong."})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
