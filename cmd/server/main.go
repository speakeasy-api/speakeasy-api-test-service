package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/speakeasy-api/speakeasy-api-test-service/internal/acceptHeaders"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/clientcredentials"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/ecommerce"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/errors"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/eventstreams"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/method"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/middleware"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/pagination"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/readonlywriteonly"
	"github.com/speakeasy-api/speakeasy-api-test-service/internal/reflect"
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
	r.HandleFunc("/oauth2/token", auth.HandleOAuth2InspectToken).Methods(http.MethodGet)
	r.HandleFunc("/oauth2/token", auth.HandleOAuth2).Methods(http.MethodPost)
	r.HandleFunc("/auth", auth.HandleAuth).Methods(http.MethodPost)
	r.HandleFunc("/auth/customsecurity/{customSchemeType}", auth.HandleCustomAuth).Methods(http.MethodGet)
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}).Methods(http.MethodGet)
	r.HandleFunc("/requestbody", requestbody.HandleRequestBody).Methods(http.MethodPost)
	r.HandleFunc("/vendorjson", responseHeaders.HandleVendorJsonResponseHeaders).Methods(http.MethodGet)
	r.HandleFunc("/pagination/limitoffset/page", pagination.HandleLimitOffsetPage).Methods(http.MethodGet, http.MethodPut)
	r.HandleFunc("/pagination/limitoffset/deep_outputs/page", pagination.HandleLimitOffsetDeepOutputsPage).Methods(http.MethodGet, http.MethodPut)
	r.HandleFunc("/pagination/limitoffset/offset", pagination.HandleLimitOffsetOffset).Methods(http.MethodGet, http.MethodPut)
	r.HandleFunc("/pagination/cursor", pagination.HandleCursor).Methods(http.MethodGet, http.MethodPut)
	r.HandleFunc("/pagination/url", pagination.HandleURL).Methods(http.MethodGet)
	r.HandleFunc("/pagination/cursor_non_numeric", pagination.HandleNonNumericCursor).Methods(http.MethodGet)
	r.HandleFunc("/retries", retries.HandleRetries).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/retries/after", retries.HandleRetries).Methods(http.MethodGet)
	r.HandleFunc("/errors/{status_code}", errors.HandleErrors).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/optional", acceptHeaders.HandleAcceptHeaderMultiplexing).Methods(http.MethodGet)
	r.HandleFunc("/readonlyorwriteonly", readonlywriteonly.HandleReadOrWrite).Methods(http.MethodPost)
	r.HandleFunc("/readonlyandwriteonly", readonlywriteonly.HandleReadAndWrite).Methods(http.MethodPost)
	r.HandleFunc("/writeonlyoutput", readonlywriteonly.HandleWriteOnlyOutput).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/json", eventstreams.HandleEventStreamJSON).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/text", eventstreams.HandleEventStreamText).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/multiline", eventstreams.HandleEventStreamMultiLine).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/rich", eventstreams.HandleEventStreamRich).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/chat", eventstreams.HandleEventStreamChat).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/chat-chunked", eventstreams.HandleEventStreamChat).Methods(http.MethodPost)
	r.HandleFunc("/eventstreams/differentdataschemas", eventstreams.HandleEventStreamDifferentDataSchemas).Methods(http.MethodPost)
	r.HandleFunc("/clientcredentials/token", clientcredentials.HandleTokenRequest).Methods(http.MethodPost)
	r.HandleFunc("/clientcredentials/authenticatedrequest", clientcredentials.HandleAuthenticatedRequest).Methods(http.MethodPost)
	r.HandleFunc("/clientcredentials/alt/token", clientcredentials.HandleTokenRequest).Methods(http.MethodPost)
	r.HandleFunc("/clientcredentials/alt/authenticatedrequest", clientcredentials.HandleAuthenticatedRequest).Methods(http.MethodPost)
	r.HandleFunc("/reflect", reflect.HandleReflect).Methods(http.MethodPost)
	r.HandleFunc("/method/delete", method.HandleDelete).Methods(http.MethodDelete)
	r.HandleFunc("/method/get", method.HandleGet).Methods(http.MethodGet)
	r.HandleFunc("/method/head", method.HandleHead).Methods(http.MethodHead)
	r.HandleFunc("/method/options", method.HandleOptions).Methods(http.MethodOptions)
	r.HandleFunc("/method/patch", method.HandlePatch).Methods(http.MethodPatch)
	r.HandleFunc("/method/post", method.HandlePost).Methods(http.MethodPost)
	r.HandleFunc("/method/put", method.HandlePut).Methods(http.MethodPut)
	r.HandleFunc("/method/trace", method.HandleTrace).Methods(http.MethodTrace)

	oauth2router := r.NewRoute().Subrouter()
	oauth2router.Use(middleware.OAuth2)
	oauth2router.HandleFunc("/ecommerce/products", ecommerce.HandleListProducts).Methods(http.MethodGet)
	oauth2router.HandleFunc("/ecommerce/products", ecommerce.HandleCreateProduct).Methods(http.MethodPost)
	oauth2router.HandleFunc("/ecommerce/products/{id}", ecommerce.HandleFetchProduct).Methods(http.MethodGet)
	oauth2router.HandleFunc("/ecommerce/products/{id}", ecommerce.HandleDeleteProduct).Methods(http.MethodDelete)
	oauth2router.HandleFunc("/ecommerce/products/{id}/inventory", ecommerce.HandleUpdateProductStock).Methods(http.MethodPut)

	handler := middleware.Fault(r)
	handler = middleware.Teapot(handler)

	bind := ":8080"
	if bindArg != nil {
		bind = *bindArg
		if !strings.HasPrefix(bind, ":") {
			bind = ":" + bind
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go auth.StartTokenDBCompaction(ctx)

	log.Printf("Listening on %s\n", bind)
	if err := http.ListenAndServe(bind, handler); err != nil {
		log.Fatal(err)
	}
}
