package authentication

// LoginReq contains Login details
//
//swagger:parameters login
type LoginReq struct {

	// Username of user
	// required: true
	// in: body
	Username string `schema:"username,required"`

	// Password of user
	// required: true
	// in: body
	Password string `schema:"password,required"`
}
