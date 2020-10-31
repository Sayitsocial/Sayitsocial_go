package common

import (
	"encoding/json"
	"net/http"
)

// Body contains body of response
type Body struct {
	// The validation message
	//
	// Required: true
	Message string `json:"message"`
}

// SuccessResponse ...
// swagger:response SuccessResponse
type SuccessResponse struct {
	// Success message
	// in: body
	Body Body
}

// ErrorResponse ...
// swagger:response errResponse
type ErrorResponse struct {
	// The error message
	// in: body
	Body Body
}

func writeBody(message string) Body {
	return Body{
		Message: message,
	}
}

// WriteSuccess writes success to http.ResponseWriter
func WriteSuccess(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&SuccessResponse{
		Body: Body{
			Message: "success",
		},
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
		Body: Body{
			Message: message,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
