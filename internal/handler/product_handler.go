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

const defaultTimeout = 5 * time.Second

type ProductHandler struct {
	service *service.ProductService // a pointer to the ProductService struct
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{service: productService} // product_Service pointer passed in main.go
}

// This function prepares the json response without having to hardcode it every single time
func sendResponse(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	switch data.(type) {
	case map[string]string:
		// If it's a map, it contains error message, no need to use marshal
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(data)
	default:
		// It is the data needed to be sent in response, therefore marshal will be used here
		result, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to encode response."})
			return
		}
		w.WriteHeader(code)
		w.Write(result)
	}
}

// Creating a new Product in collection
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Creating context with timeout of 5 second
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel() // cancel() helps to release resources in case if the request is completed before the context is expired.

	// Extracting product from request
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	// Creating new product by passing it through to the service layer
	if err := h.service.CreateProduct(ctx, &product); err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}) // err.Error() returns the error message sent by service
		return
	}

	// setting up response
	sendResponse(w, http.StatusCreated, product)
}

// Get all the products with a particular category
func (h *ProductHandler) GetProductByCategory(w http.ResponseWriter, r *http.Request) {
	// context
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()
	w.Header().Set("Content-Type", "application/json")

	// Fetching category title from URL
	params := mux.Vars(r)
	categoryTitle := params["category"]
	if categoryTitle == "" {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Category name is required."})
		return
	}

	products, err := h.service.GetByCategory(ctx, categoryTitle)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products."})
		return
	}

	// Encoding response
	sendResponse(w, http.StatusOK, products)
}

// Get a specific product with it's ID provided in URL path
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	// Extracting ID from URL path
	params := mux.Vars(r)
	reqID := params["id"]
	if reqID == "" {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	// Passing request and context to service layer
	product, err := h.service.GetProduct(ctx, reqID)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong."})
		return
	}

	sendResponse(w, http.StatusOK, product)
}

// Get a slice of all the products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	// passing context to service layer
	products, err := h.service.ListProducts(ctx)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
		return
	}

	sendResponse(w, http.StatusOK, products)
}

// Updating a particular product with the product's ID given in URL
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	// Extracting ID from URL path
	params := mux.Vars(r)
	reqID := params["id"]
	if reqID == "" {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID."})
		return
	}

	// Extracting product from request
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request."})
		return
	}

	// Passing request to service layer
	if err := h.service.UpdateProduct(ctx, reqID, &product); err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong."})
		return
	}
	sendResponse(w, http.StatusOK, product)
}

// Delete a specific product with it's ID given in URL path
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	// Extracting ID from URL path
	params := mux.Vars(r)
	reqID := params["id"]
	if reqID == "" {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID."})
		return
	}

	// Passing request to service layer
	if err := h.service.DeleteProduct(ctx, reqID); err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong."})
		return
	}

	sendResponse(w, http.StatusNoContent, nil)
}
