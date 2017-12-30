package googleSignInIDTokenVerifier

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	DefaultRefreshDuration = time.Hour
)

var (
	// Issuers default issuers checked against during the verification process
	Issuers = []string{"accounts.google.com", "https://accounts.google.com"}
)

// Verifier can be used to modifiy the behavior of the verifier
type Verifier struct {
	LazyLoad    bool
	keys        map[string]*rsa.PublicKey
	cacheExpiry time.Time

	ticker *time.Ticker
}

// SharedInstance can be used globally, is lazy loading
var SharedInstance = NewVerifier()

// NewVerifier returns a new verifier
func NewVerifier() *Verifier {
	return &Verifier{
		LazyLoad: true,
	}
}

// LazyLoading can be chained after NewVerifier and set the loading behavior
// Disables any runing timer from PeriodicRefresh
// lazy == true: data is loaded when needed and cached (e.g. during verify requests) can make some requests longer
// lazy == false: data is not loaded automatically, you should either use PeriodicRefresh or refresh keys manually
func (v *Verifier) LazyLoading(lazy bool) *Verifier {
	v.LazyLoad = lazy
	if v.ticker != nil {
		v.ticker.Stop()
	}
	return v
}

// PeriodicRefresh sets a ticker every 'duration' refreshing keys.
// LazyLoaging is automatically disabled
// To stop the timer, call the LazyLoading method with any parameter
// Can be chained
func (v *Verifier) PeriodicRefresh(duration time.Duration) *Verifier {
	v.LazyLoading(false)
	v.ticker = time.NewTicker(duration)

	go func() {
		for range v.ticker.C {
			v.ForceRefreshCerts()
		}
		v.ticker = nil
	}()

	return v
}

// Verify verifies the given token with Google keys with the shared Verifier
func Verify(idToken string, audience ...string) error {
	return SharedInstance.Verify(idToken, audience...)
}

// Verify verifies the given token with Google keys with the shared Verifier
func (v *Verifier) Verify(idToken string, audience ...string) error {
	_, err := v.Decode(idToken, audience...)
	return err
}

// Decode verifies the given Google jwt token and returns its claims with the shared Verifier
func Decode(idToken string, audience ...string) (*GoogleClaims, error) {
	return SharedInstance.Decode(idToken, audience...)
}

// Decode verifies the given Google jwt token and returns its claims
func (v *Verifier) Decode(idToken string, audience ...string) (*GoogleClaims, error) {
	token, err := jwt.ParseWithClaims(idToken, &GoogleClaims{}, func(token *jwt.Token) (interface{}, error) {
		if kid, ok := token.Header["kid"].(string); ok {
			if v.LazyLoad {
				_, err := v.RefreshCerts()
				if err != nil {
					return nil, err
				}
			}
			fmt.Println(v.keys)
			if key, ok := v.keys[kid]; ok {
				return key, nil
			}
		}
		return nil, ErrPublicKeyNotFound
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*GoogleClaims)

	var ok bool
	for _, issuer := range Issuers {
		ok = claims.VerifyIssuer(issuer, true)
		if ok {
			break
		}
	}
	if !ok {
		return nil, ErrInvalidIssuer
	}

	ok = false
	for _, aud := range audience {
		ok = claims.VerifyAudience(aud, true)
		if ok {
			break
		}
	}
	if !ok {
		return nil, ErrInvalidAudience
	}

	return claims, nil
}
