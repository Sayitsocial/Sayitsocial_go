package authentication

// LoginReq contains Login details
//
//swagger:parameters login JWTLogin
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

// JWTResp response for JWT logins
//
//swagger:response JWTLoginResp
type JWTResp struct {
	// Token provided after successful login
	Token string `json:"token"`
}

// JWTRefreshReq contains Token for refresh
//
//swagger:parameters login JWTRefresh
type JWTRefreshReq struct {
	// in: query
	// required: true
	// Token provided after successful login
	Token string `json:"token"`
}
