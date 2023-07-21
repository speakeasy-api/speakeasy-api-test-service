package acceptHeaders

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/speakeasy-api/speakeasy-api-test-service/internal/utils"
)

func HandleAcceptHeaderMultiplexing(w http.ResponseWriter, r *http.Request) {
	var obj interface{}
	if r.Header["Accept"] == "application/json" {
		err := json.Unmarshal([]byte("{\"Obj1\":\"obj1\"}"), &obj)
		if err != nil {
			utils.HandleError(w, err)
			return
		}
	
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
		if err := json.NewEncoder(w).Encode(obj); err != nil {
			utils.HandleError(w, err)
		}
	} else if r.Header["Accept"] == "application/xml" {
		err := json.Unmarshal([]byte("{\"Obj2\":\"obj2\"}"), &obj)
		if err != nil {
			utils.HandleError(w, err)
			return
		}
	
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	
		if err := xml.NewEncoder(w).Encode(obj); err != nil {
			utils.HandleError(w, err)
		}
	}

}
