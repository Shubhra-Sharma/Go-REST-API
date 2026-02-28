package handler

import (
    "context"
    "encoding/json"
    "github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
    "github.com/Shubhra-Sharma/Go-REST-API/internal/service"
    "net/http"
    "time"
    "github.com/gorilla/mux"
)

type ProductHandler struct {
    service *service.ProductService // a pointer to the ProductService interface
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
    return &ProductHandler{service: service} // service pointer passed in main.go after creating ProductService interface
}


// Creating a new Product in collection
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    // Creating context with timeout of 5 second
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    // Extracting product from request
    var product domain.Product
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
       w.WriteHeader(http.StatusBadRequest)
       json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request."})
	   return
    }
    
    // Creating new product by passing it through to the service layer
    if err := h.service.CreateProduct(ctx, &product); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong."})
        return
    }
    
    // setting up response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated) // How to set status after encoding?
    json.NewEncoder(w).Encode(product) 
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
    
    //Passing request to service layer
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