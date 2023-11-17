package auth

import (
	"encoding/json"
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
