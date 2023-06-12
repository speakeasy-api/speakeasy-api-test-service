package main

import (
	"log"
	"net/http"

	"github.com/speakeasy-api/speakeasy-auth-test-service/internal/pagination"
	"github.com/speakeasy-api/speakeasy-auth-test-service/internal/responseHeaders"
	"github.com/speakeasy-api/speakeasy-auth-test-service/internal/retries"

	"github.com/gorilla/mux"
	"github.com/speakeasy-api/speakeasy-auth-test-service/internal/auth"
	"github.com/speakeasy-api/speakeasy-auth-test-service/internal/requestbody"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}).Methods(http.MethodGet)
	r.HandleFunc("/auth", auth.HandleAuth).Methods(http.MethodPost)
	r.HandleFunc("/requestbody", requestbody.HandleRequestBody).Methods(http.MethodPost)
	r.HandleFunc("/vendorjson", responseHeaders.HandleVendorJsonResponseHeaders).Methods(http.MethodGet)
	r.HandleFunc("/pagination/limitoffset/page", pagination.HandleLimitOffsetPage).Methods(http.MethodGet, http.MethodPut)
	r.HandleFunc("/pagination/limitoffset/offset", pagination.HandleLimitOffsetOffset).Methods(http.MethodGet, http.MethodPut)
	r.HandleFunc("/pagination/cursor", pagination.HandleCursor).Methods(http.MethodGet, http.MethodPut)
	r.HandleFunc("/retries", retries.HandleRetries).Methods(http.MethodGet)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
