package authentication

type Body struct {
	// The validation message
	//
	// Required: true
	// Example: Expected type int
	Message string `json:"message"`
}

// Invalid Credentials
// swagger:response successResponse
type invalidCredentialsError struct {
	// The error message
	// in: body
	Body Body
}
