package method

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/speakeasy-api/speakeasy-api-test-service/internal/utils"
)

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		ID string `json:"id"`
	}
	type ResponseBody struct {
		Status string `json:"status"`
	}

	var requestBody RequestBody

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.HandleError(w, err)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseBody := ResponseBody{
		Status: "OK",
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		utils.HandleError(w, err)
		return
	}
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		ID string `json:"id"`
	}
	type ResponseBody struct {
		Status string `json:"status"`
	}

	var requestBody RequestBody

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.HandleError(w, err)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseBody := ResponseBody{
		Status: "OK",
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		utils.HandleError(w, err)
		return
	}
}

func HandleHead(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		ID string `json:"id"`
	}

	var requestBody RequestBody

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.HandleError(w, err)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func HandleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "OPTIONS")
	w.WriteHeader(http.StatusOK)
}

func HandlePatch(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		ID string `json:"id"`
	}
	type ResponseBody struct {
		Status string `json:"status"`
	}

	var requestBody RequestBody

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.HandleError(w, err)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseBody := ResponseBody{
		Status: "OK",
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		utils.HandleError(w, err)
		return
	}
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		ID string `json:"id"`
	}
	type ResponseBody struct {
		Status string `json:"status"`
	}

	var requestBody RequestBody

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.HandleError(w, err)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseBody := ResponseBody{
		Status: "OK",
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		utils.HandleError(w, err)
		return
	}
}

func HandlePut(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		ID string `json:"id"`
	}
	type ResponseBody struct {
		Status string `json:"status"`
	}

	var requestBody RequestBody

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.HandleError(w, err)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseBody := ResponseBody{
		Status: "OK",
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		utils.HandleError(w, err)
		return
	}
}

func HandleTrace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "message/http")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("TRACE /method/trace HTTP/1.1\r\nHost: example.com\r\n"))
}
