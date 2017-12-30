package googleSignInIDTokenVerifier

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	// GoogleCertsURL url exposed by Google with their current RSA keys
	GoogleCertsURL = "https://www.googleapis.com/oauth2/v3/certs"

	defaultCacheDuration = time.Hour * 2
)

type gKey struct {
	Kty string
	Alg string
	Use string
	Kid string
	N   string
	E   string
}

func (v *Verifier) refreshCerts() error {
	resp, err := http.Get(GoogleCertsURL)
	if err != nil {
		return err
	}

	v.cacheExpiry = time.Now().Add(defaultCacheDuration)
	if cache := resp.Header.Get("cache-control"); cache != "" {
		r := regexp.MustCompile("max-age=([0-9]*)")
		matches := r.FindAllStringSubmatch(cache, -1)

		if len(matches) > 0 && len(matches[0]) > 1 {
			var offset int64
			offset, err = strconv.ParseInt(matches[0][1], 10, 64)
			if err == nil {
				v.cacheExpiry = time.Now().Add(time.Duration(offset) * time.Second)
			}
		}
	}

	var dest map[string][]gKey
	err = json.NewDecoder(resp.Body).Decode(&dest)
	if err != nil {
		return err
	}

	if len(dest["keys"]) < 1 {
		return fmt.Errorf("No keys were found when fetching %s", GoogleCertsURL)
	}

	keys := map[string]*rsa.PublicKey{}
	for _, gkey := range dest["keys"] {
		if gkey.Kty == "RSA" && gkey.Use == "sig" {
			n, err := base64.RawURLEncoding.DecodeString(gkey.N)
			if err != nil {
				return err
			}

			e, err := base64.RawURLEncoding.DecodeString(gkey.E)
			if err != nil {
				return err
			}

			keys[gkey.Kid] = &rsa.PublicKey{
				N: big.NewInt(0).SetBytes([]byte(n)),
				E: int(big.NewInt(0).SetBytes([]byte(e)).Int64()),
			}
		}
	}

	v.keys = keys
	return nil
}

// RefreshCerts refreshes current certificates if keys are expired
// Returns a boolean whether it hit cache or not
func (v *Verifier) RefreshCerts() (bool, error) {
	if time.Now().Before(v.cacheExpiry) {
		return true, nil
	}
	return false, v.refreshCerts()
}

// RefreshCerts refreshes current certificates if keys are expired
// Returns a boolean whether it hit cache or not
func RefreshCerts() (bool, error) {
	return SharedInstance.RefreshCerts()
}

// ForceRefreshCerts forcefully refreshes certificates
func (v *Verifier) ForceRefreshCerts() error {
	return v.refreshCerts()
}

// ForceRefreshCerts forcefully refreshes certificates
func ForceRefreshCerts() error {
	return SharedInstance.ForceRefreshCerts()
}
