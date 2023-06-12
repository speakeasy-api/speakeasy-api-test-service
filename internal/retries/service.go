package retries

import (
	"net/http"
	"strconv"
)

var callCounts = map[string]int{}

func HandleRetries(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("request-id")
	numRetriesStr := r.URL.Query().Get("num-retries")

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

	_, ok := callCounts[requestID]
	if !ok {
		callCounts[requestID] = 0
	}
	callCounts[requestID]++

	if callCounts[requestID] < numRetries {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("request failed please retry"))
		return
	}

	delete(callCounts, requestID)
	w.WriteHeader(http.StatusOK)
}