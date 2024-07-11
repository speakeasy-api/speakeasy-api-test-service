package models

type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Type    string `json:"type"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type ErrorType1Response struct {
	Error string `json:"error"`
}

type ErrorType2Response struct {
	Error ErrorMessage `json:"error"`
}

type TaggedError1Response struct {
	Tag   string `json:"tag"`
	Error string `json:"error"`
}

type TaggedError2Response struct {
	Tag   string       `json:"tag"`
	Error ErrorMessage `json:"error"`
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
