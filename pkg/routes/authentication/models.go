package authentication

// LoginReq contains Login details
//
//swagger:parameters login
type LoginReq struct {

	// Username of user
	// required: true
	// in: query
	Username string `schema:"username,required" json:"username"`

	// Password of user
	// required: true
	// in: query
	Password string `schema:"password,required" json:"password"`
}
