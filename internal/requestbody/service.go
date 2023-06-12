package requestbody

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/speakeasy-api/speakeasy-api-test-service/internal/utils"
)

func HandleRequestBody(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	var req interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(req); err != nil {
		utils.HandleError(w, err)
	}
}
