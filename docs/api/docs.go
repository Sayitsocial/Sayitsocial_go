// Package api classification SayItSocial.
//
// Documentation of SayItSocialAPI.
//
//     Schemes: http
//     BasePath:
//     Version: 1.0.0
//     Host: sayitsocial-dev.ddns.net:8000
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - cookieAuth
//     - JWTAuth
//
//     SecurityDefinitions:
//       cookieAuth:
//         type: apiKey
//         in: cookie
//         name: SESSIONID
//         description: Doesnt work from swagger, head over to /login
//       JWTAuth:
//         type: apiKey
//         in: header
//         name: Authorization
//         description: Get token from /auth/jwt-login
//
// swagger:meta
package api
