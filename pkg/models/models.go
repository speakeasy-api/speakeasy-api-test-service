package models

type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type HeaderAuth struct {
	HeaderName    string `json:"headerName"`
	ExpectedValue string `json:"expectedValue"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRequest struct {
	HeaderAuth []HeaderAuth `json:"headerAuth,omitempty"`
	BasicAuth  *BasicAuth   `json:"basicAuth,omitempty"`
}
