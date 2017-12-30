package googleSignInIDTokenVerifier

import "github.com/dgrijalva/jwt-go"

// GoogleClaims exposes all the fields set by Google about your token
type GoogleClaims struct {
	jwt.StandardClaims
	AuthorizedParty string `json:"azp,omitempty"`

	Name       string `json:"name,omitempty"`
	GivenName  string `json:"given_name,omitempty"`
	FamilyName string `json:"family_name,omitempty"`

	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	Picture       string `json:"picture,omitempty"`
	Locale        string `json:"locale,omitempty"`
}
