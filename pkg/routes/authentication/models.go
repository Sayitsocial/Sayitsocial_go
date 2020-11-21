package authentication

// LoginReq contains Login details
//
//swagger:model
type LoginReq struct {
	// Username of user
	// required: true
	Username string `schema:"username,required" json:"username"`

	// Password of user
	// required: true
	Password string `schema:"password,required" json:"password"`
}

//swagger:parameters login JWTLogin
type loginModel struct {
	// in: body
	Login LoginReq
}

// JWTResp response for JWT logins
//
//swagger:response JWTLoginResp
type JWTResp struct {
	// Token provided after successful login
	Token string `json:"token"`
}

// JWTRefreshReq contains Token for refresh
// swagger:model
type JWTRefreshReq struct {
	// Token provided after successful login
	// required: true
	Token string `json:"token"`
}

//swagger:parameters JWTRefresh
type jwtModel struct {
	// in: body
	JWT JWTRefreshReq
}
