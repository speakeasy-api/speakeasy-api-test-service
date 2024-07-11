package errors

import (
	"encoding/json"
	"fmt"
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

	if err := json.NewEncoder(w).Encode(models.Error{
		Code:    statusCode,
		Message: "an error occurred",
		Type:    "internal",
	}); err != nil {
		utils.HandleError(w, err)
		return
	}
}

func HandleUnionOfErrors(w http.ResponseWriter, r *http.Request) {
	errorType := r.FormValue("errorType")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	var res interface{}
	switch errorType {
	case "type1":
		res = models.ErrorType1Response{
			Error: "Error1",
		}
	case "type2":
		res = models.ErrorType2Response{
			Error: models.ErrorMessage{
				Message: "Error2",
			},
		}
	default:
		utils.HandleError(w, fmt.Errorf("unknown error type: \"type1\" or \"type2\" expected"))
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
}

func HandleDiscriminatedUnionOfErrors(w http.ResponseWriter, r *http.Request) {
	errorTag := r.FormValue("errorTag")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	var res interface{}
	switch errorTag {
	case "tag1":
		res = models.TaggedError1Response{
			Tag:   "tag1",
			Error: "Error1",
		}
	case "tag2":
		res = models.TaggedError2Response{
			Tag: "tag2",
			Error: models.ErrorMessage{
				Message: "Error2",
			},
		}
	default:
		utils.HandleError(w, fmt.Errorf("unknown error tag: \"tag1\" or \"tag2\" expected"))
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
}
