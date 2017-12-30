package googleSignInIDTokenVerifier

import (
	"testing"
	"time"
)

func TestRefreshCerts(t *testing.T) {
	err := ForceRefreshCerts()
	if err != nil {
		t.Error(err)
	}

	cached, err := RefreshCerts()
	if err != nil {
		t.Error(err)
	}
	if !cached {
		t.Errorf("Did not hit cache")
	}

	SharedInstance.cacheExpiry = time.Now().Truncate(time.Second * 5)
	cached, err = RefreshCerts()
	if err != nil {
		t.Error(err)
	}
	if cached {
		t.Errorf("Expired keys were returned")
	}
}
