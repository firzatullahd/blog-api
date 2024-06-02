package response

import (
	"encoding/json"
	"net/http"
)

type HTTPResponse struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func SetHTTPResponse(w http.ResponseWriter, statusCode int, msg string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(HTTPResponse{
		StatusCode: statusCode,
		Message:    msg,
		Data:       data,
	})
}
