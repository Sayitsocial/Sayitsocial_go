package common

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse ...
// swagger:response successResponse
type SuccessResponse struct {
	// Success message
	// in: body
	// Required: true
	Message string `json:"message"`
}

// ErrorResponse ...
// swagger:response errResponse
type ErrorResponse struct {
	// The error message
	// in: body
	// Required: true
	Message string `json:"message"`
}

// WriteSuccess writes success to http.ResponseWriter
func WriteSuccess(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&SuccessResponse{
		Message: "success",
	})
	if err != nil {
		return err
	}
	return nil
}

// WriteError writes error to http.ResponseWriter
func WriteError(message string, statusCode int, w http.ResponseWriter) error {
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&ErrorResponse{
		Message: message,
	})
	if err != nil {
		return err
	}
	return nil
}
