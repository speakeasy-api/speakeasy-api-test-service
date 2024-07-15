package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/utils"

	"github.com/speakeasy-api/speakeasy-api-test-service/pkg/models"
)

func HandleErrors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statusCode, ok := vars["status_code"]
	if !ok {
		utils.HandleError(w, fmt.Errorf("status_code is required"))
		return
	}

	statusCodeInt, err := strconv.Atoi(statusCode)
	if err != nil {
		utils.HandleError(w, fmt.Errorf("status_code must be an integer"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeInt)

	var res interface{}
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		if err := json.Unmarshal(body, &res); err != nil {
			utils.HandleError(w, err)
			return
		}
	} else {
		res = models.Error{
			Code:    statusCode,
			Message: "an error occurred",
			Type:    "internal",
		}
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.HandleError(w, err)
		return
	}
}
