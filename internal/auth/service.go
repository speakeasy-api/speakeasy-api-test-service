package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/speakeasy-api/speakeasy-api-test-service/internal/utils"

	"github.com/speakeasy-api/speakeasy-api-test-service/pkg/models"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	var req models.AuthRequest
	if err := json.Unmarshal(body, &req); err != nil {
		utils.HandleError(w, err)
		return
	}

	if err := checkAuth(req, r); err != nil {
		utils.HandleError(w, err)
		return
	}
}

func HandleCustomAuth(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/auth/customsecurity/customSchemeAppId":
		appID := r.Header.Get("X-Security-App-Id")
		if appID != "testAppID" {
			utils.HandleError(w, fmt.Errorf("invalid app id: %w", authError))
			return
		}
		secret := r.Header.Get("X-Security-Secret")
		if secret != "testSecret" {
			utils.HandleError(w, fmt.Errorf("invalid secret: %w", authError))
			return
		}
	default:
		utils.HandleError(w, fmt.Errorf("invalid path"))
		return
	}
}
