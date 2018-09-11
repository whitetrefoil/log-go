package server

import (
	"encoding/json"
	"net/http"
)

type responseWrapper struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func doResponseWithMessage(w http.ResponseWriter, code int, message string, body interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(&responseWrapper{
		Code:    code,
		Message: message,
		Data:    body,
	})
}

func doResponse(w http.ResponseWriter, code int, body interface{}) error {
	return doResponseWithMessage(w, code, http.StatusText(code), body)
}
