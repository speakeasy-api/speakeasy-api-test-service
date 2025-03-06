package clientcredentials

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"sync"
)

var state = sync.Map{}

const (
	firstAccessToken  = "super-duper-access-token"
	secondAccessToken = "second-super-duper-access-token"
)

func handleBasicAuth(authHeader string) (clientID, clientSecret string, ok bool) {
	// Remove "Basic " prefix (case-insensitive)
	if !strings.HasPrefix(strings.ToLower(authHeader), "basic ") {
		return "", "", false
	}
	encodedCreds := strings.TrimSpace(authHeader[6:])

	// Decode base64
	decodedCreds, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encodedCreds))
	if err != nil {
		return "", "", false
	}

	// Split into username:password
	creds := strings.SplitN(string(decodedCreds), ":", 2)
	if len(creds) != 2 {
		return "", "", false
	}

	return creds[0], creds[1], true
}


func HandleTokenRequest(w http.ResponseWriter, r *http.Request) {
	var clientID, clientSecret string
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	// Check for Basic Auth header
	if authHeader := r.Header.Get("Authorization"); strings.HasPrefix(strings.ToLower(authHeader), "basic ") {
		var ok bool
		clientID, clientSecret, ok = handleBasicAuth(authHeader)
		if !ok {
			http.Error(w, "invalid_client", http.StatusUnauthorized)
			return
		}
	} else {
		clientID = r.Form.Get("client_id")
		clientSecret = r.Form.Get("client_secret")
	}
	grant_type := r.Form.Get("grant_type")
	if grant_type != "client_credentials" {
		http.Error(w, "invalid_grant", http.StatusBadRequest)
		return
	}



	if clientID == "" || clientSecret == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}
	if clientID != "speakeasy-sdks" || !strings.HasPrefix(clientSecret, "supersecret-") {
		http.Error(w, "invalid_client", http.StatusUnauthorized)
		return
	}

	scopes := strings.Split(r.Form.Get("scope"), " ")
	if len(scopes) == 0 {
		http.Error(w, "invalid_scope", http.StatusBadRequest)
		return
	}

	if !slices.Contains(scopes, "read") && !slices.Contains(scopes, "write") {
		http.Error(w, "invalid_scope", http.StatusBadRequest)
		return
	}

	accessToken := firstAccessToken

	_, ok := state.Load(clientSecret)
	if !ok {
		state.Store(clientSecret, true)
	} else {
		accessToken = secondAccessToken
	}

	w.Header().Set("Content-Type", "application/json")

	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	response := tokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   0,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "server_error", http.StatusInternalServerError)
		return
	}
}

func HandleAuthenticatedRequest(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	if accessToken != firstAccessToken && accessToken != secondAccessToken {
		http.Error(w, "invalid_token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
