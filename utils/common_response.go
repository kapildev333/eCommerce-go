package utils

type Response struct {
	StatusCode int         `json:"status_code"`
	Error      bool        `json:"error"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}
