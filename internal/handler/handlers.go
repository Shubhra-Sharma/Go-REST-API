package handlers
import(
	"encoding/json"
	"net/http"
	"github.com/Shubhra-Sharma/Go-REST-API/internal/domain"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.ProductMap)
	
}
//r = incoming http request
//encoding/json
//This package handles converting between Go data structures and JSON text. It's used in two directions here:
//Encoding (Go → JSON, for responses)
//Decoding (JSON → Go, for request bodies)
// NewDecoder(r.Body) reads the raw JSON from the incoming request body, and .Decode(&newProduct) parses it into your Go struct. The `&` passes a pointer so the struct is actually populated.


//Get a particular product
func ProductHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json") // this tells client that the body sent as response is in json

	reqID, err := utils.FetchId(w,r)   // id is a string and is converted to an integer for comparison
	//error handling
	if err != nil {
		return
	}

	if _,present := model.ProductMap[reqID]; present{
		    w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(model.ProductMap[reqID])
			return
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
}

//CREATE
func AddProductHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
    var newProduct model.Product
	err:= json.NewDecoder(r.Body).Decode(&newProduct)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
       json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request."})
	   return
	}
    
	if newProduct.Name == "" || newProduct.Price <= 0 || newProduct.Quantity < 0 {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]string{"error": "Name is required, price must be > 0 and quantity must be >= 0"})
    return
    }

	newID := len(model.ProductMap) + 1
	newProduct.ID = newID
	model.ProductMap[newProduct.ID] = newProduct
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

//PUT
func UpdateProductHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	reqID, err := utils.FetchId(w,r)
	if err != nil {
		return
	}

	//Extracting the data to be updated
	var updatedProduct model.Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request."})
		return
	}
    
	if updatedProduct.Name == "" || updatedProduct.Price <= 0 || updatedProduct.Quantity < 0 {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]string{"error": "Name is required, price must be > 0 and quantity must be >= 0"})
    return
    }
	
	 _, present := model.ProductMap[reqID]

	 if present {
	updatedProduct.ID = reqID
    model.ProductMap[reqID] = updatedProduct
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedProduct)
    return
    }

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string {"error": "Product not found"})
}


//DELETE
func DeleteProductHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	reqID,err := utils.FetchId(w,r)
	if err!=nil {
		return
	}

	if _,present := model.ProductMap[reqID]; present{
		delete(model.ProductMap,reqID)
		w.WriteHeader(http.StatusNoContent) //The status code for http.StatusNoContent is 204. It means the request was successful but there's nothing to send back in the response body
		return
	}

	w.WriteHeader(http.StatusNotFound)  
	json.NewEncoder(w).Encode(map[string]string {"error": "Product not found"})
}