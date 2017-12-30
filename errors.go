package googleSignInIDTokenVerifier

import "errors"

var ErrPublicKeyNotFound = errors.New("Public key for given key id not found")

var ErrInvalidIssuer = errors.New("Issuer not provided or invalid")

var ErrInvalidAudience = errors.New("Audience not provided or invalid")
