package utils
import(
	"strconv"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)
func FetchId(w http.ResponseWriter,r *http.Request) (int,error){
	params := mux.Vars(r)
	reqID, err := strconv.Atoi(params["id"])   // id is a string and is converted to an integer for comparison
	//error handling
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid product ID"})
		return -1,err
	}
	return reqID,nil
}