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

type ProductHandler struct {
	service         *service.ProductService // a pointer to the ProductService struct
	categoryService *service.ProductCategoryService
}

func NewProductHandler(product_Service *service.ProductService, category_service *service.ProductCategoryService) *ProductHandler {
	return &ProductHandler{
		service:         product_Service,
		categoryService: category_service,
	} // product_Service pointer and category_service pointer passed in main.go
}

// Creating a new Product in collection
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Creating context with timeout of 5 second
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel() // cancel() helps to release resources in case if the request is completed before the context is expired.

	// Extracting product from request
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request."})
		return
	}
	if product.Category == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please choose a category for your product"})
		return
	}

	// First fetch the categoryID for the respective category
	categoryID, err := h.categoryService.GetCategoryID(ctx, product.Category)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "category not found"})
		return
	}
	product.CategoryID = categoryID
	// Creating new product by passing it through to the service layer
	if err := h.service.CreateProduct(ctx, &product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong."})
		return
	}

	// setting up response
	result, err := json.Marshal(product) // Marshal encodes GO's data structures into json format and returns a slice of bytes
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to encode response"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

// Get all the products with a particular category
func (h *ProductHandler) GetProductByCategory(w http.ResponseWriter, r *http.Request) {
	// context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	w.Header().Set("Content-Type", "application/json")

	// Fetching category title from URL
	categoryTitle := r.URL.Query().Get("category")
	if categoryTitle == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "category name is required"})
		return
	}

	// Fetching categoryID from category collection so that filtering can be done using categoryID
	categoryID, err := h.categoryService.GetCategoryID(ctx, categoryTitle)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "category not found"})
		return
	}

	// After fetchinf categoryID, passing context to service layer
	products, err := h.service.GetByCategory(ctx, categoryID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch products"})
		return
	}

	// Encoding response
	result, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to encode response"})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Get a specific product with it's ID provided in URL path
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
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

	// Passing request and context to service layer
	product, err := h.service.GetProduct(ctx, reqID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong."})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Get a slice of all the products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// passing context to service layer
	products, err := h.service.ListProducts(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong."})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// Updating a particular product with the product's ID given in URL
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	// Extracting product from request
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request."})
		return
	}

	// Passing request to service layer
	if err := h.service.UpdateProduct(ctx, reqID, &product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong."})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Delete a specific product with it's ID given in URL path
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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
	if err := h.service.DeleteProduct(ctx, reqID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong."})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
