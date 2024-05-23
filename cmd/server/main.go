package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/speakeasy-api/speakeasy-api-test-service/internal/acceptHeaders"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/clientcredentials"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/errors"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/eventstreams"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/middleware"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/pagination"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/readonlywriteonly"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/responseHeaders"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/retries"

	"github.com/gorilla/mux"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/auth"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/requestbody"
)

var bindArg = flag.String("b", ":8080", "Bind address")

func main() {
	flag.Parse()

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
	r.HandleFunc("/pagination/url", pagination.HandleURL).Methods(http.MethodGet)
	r.HandleFunc("/pagination/cursor_non_numeric", pagination.HandleNonNumericCursor).Methods(http.MethodGet)
	r.HandleFunc("/retries", retries.HandleRetries).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/retries/after", retries.HandleRetries).Methods(http.MethodGet)
	r.HandleFunc("/errors/{status_code}", errors.HandleErrors).Methods(http.MethodGet)
	r.HandleFunc("/optional", acceptHeaders.HandleAcceptHeaderMultiplexing).Methods(http.MethodGet)
	r.HandleFunc("/readonlyorwriteonly", readonlywriteonly.HandleReadOrWrite).Methods(http.MethodPost)
	r.HandleFunc("/readonlyandwriteonly", readonlywriteonly.HandleReadAndWrite).Methods(http.MethodPost)
	r.HandleFunc("/writeonlyoutput", readonlywriteonly.HandleWriteOnlyOutput).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/json", eventstreams.HandleEventStreamJSON).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/text", eventstreams.HandleEventStreamText).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/multiline", eventstreams.HandleEventStreamMultiLine).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/rich", eventstreams.HandleEventStreamRich).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/chat", eventstreams.HandleEventStreamChat).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/differentdataschemas", eventstreams.HandleEventStreamDifferentDataSchemas).Methods(http.MethodPost)
	r.HandleFunc("/clientcredentials/token", clientcredentials.HandleTokenRequest).Methods(http.MethodPost)
	r.HandleFunc("/clientcredentials/authenticatedrequest", clientcredentials.HandleAuthenticatedRequest).Methods(http.MethodPost)

	handler := middleware.Fault(r)

	bind := ":8080"
	if bindArg != nil {
		bind = *bindArg
	}

	log.Printf("Listening on %s\n", bind)
	if err := http.ListenAndServe(bind, handler); err != nil {
		log.Fatal(err)
	}
}
