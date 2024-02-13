package retries

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var callCounts = map[string]int{}

type retriesResponse struct {
	Retries int `json:"retries"`
}

func HandleRetries(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("request-id")
	numRetriesStr := r.URL.Query().Get("num-retries")
	includeHeaderTimeout := r.URL.Query().Get("included-header-timeout")

	numRetries := 3
	if numRetriesStr != "" {
		var err error
		numRetries, err = strconv.Atoi(numRetriesStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("num-retries must be an integer"))
			return
		}
	}

	var headerTimeout int
	if includeHeaderTimeout != "" {
		var err error
		headerTimeout, err = strconv.Atoi(includeHeaderTimeout)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("included-header-timeout must be an integer"))
			return
		}
	}

	if requestID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("request-id is required"))
		return
	}

	_, ok := callCounts[requestID]
	if !ok {
		callCounts[requestID] = 0
	}
	callCounts[requestID]++

	if callCounts[requestID] < numRetries {
		if headerTimeout > 0 {
			w.Header().Set("Retry-After", strconv.Itoa(headerTimeout))
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("request failed please retry"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(retriesResponse{
		Retries: callCounts[requestID],
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failed to marshal response"))
		return
	}
	_, _ = w.Write(data)

	delete(callCounts, requestID)
}
