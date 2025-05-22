package errors

import (
	"bytes"
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

// Returns one of an errors union with a 400 Bad Request status, based on the
// provided tag name in the path parameters with these response schemas:
//
//	taggedError1:
//	  type: object
//	  properties:
//	    tag:
//	      type: string
//	      enum: [tag1]
//	    error:
//	      type: string
//	  required:
//	    - tag
//	    - error
//	taggedError2:
//	  type: object
//	    properties:
//	      tag:
//	        type: string
//	        const: tag2
//	      error:
//	        type: object
//	        properties:
//	          message:
//	            type: string
//	        required:
//	          - message
//	    required:
//	      - tag
//	      - error
func HandleErrorsUnion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag, ok := vars["tag"]

	if !ok {
		utils.HandleError(w, fmt.Errorf("tag path parameter is required"))
		return
	}

	if tag != "tag1" && tag != "tag2" {
		utils.HandleError(w, fmt.Errorf("tag path parameter must be either tag1 or tag2"))
		return
	}

	var errorRes any
	var responseBody bytes.Buffer

	if tag == "tag1" {
		errorRes = struct {
			Tag   string `json:"tag"`
			Error string `json:"error"`
		}{
			Tag:   "tag1",
			Error: "intentional tag1 error",
		}
	} else {
		errorRes = struct {
			Tag   string `json:"tag"`
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}{
			Tag: "tag2",
			Error: struct {
				Message string `json:"message"`
			}{
				Message: "intentional tag2 error",
			},
		}
	}

	if err := json.NewEncoder(&responseBody).Encode(errorRes); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(responseBody.Bytes())
}
