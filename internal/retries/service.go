package retries

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

var (
	callCounts      = map[string]int{}
	callCountsMutex sync.Mutex
)

type retriesResponse struct {
	Retries int `json:"retries"`
}

func HandleRetries(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("request-id")
	numRetriesStr := r.URL.Query().Get("num-retries")
	retryAfterVal := r.URL.Query().Get("retry-after-val")

	retryAfter := 0
	if retryAfterVal != "" {
		var err error
		retryAfter, err = strconv.Atoi(retryAfterVal)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("retry-after-val must be an integer"))
			return
		}
	}

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

	if requestID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("request-id is required"))
		return
	}

	callCountsMutex.Lock()
	_, ok := callCounts[requestID]
	if !ok {
		callCounts[requestID] = 0
	}
	callCounts[requestID]++
	callCountsMutex.Unlock()

	if callCounts[requestID] < numRetries {
		if retryAfter > 0 {
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
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

	callCountsMutex.Lock()
	delete(callCounts, requestID)
	callCountsMutex.Unlock()
}
