package googleSignInIDTokenVerifier

import (
	"testing"
	"time"
)

const (
	ValidToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjYzYjRiNzRlMDQ5OGI5MjI2NTUxOGExNjc4MWFmZGI4ZDRlZjE1ZTMifQ.eyJhenAiOiI1NDQzMDQxODc1OTYtMGxoYTRmZWwzZjY4N2szdTdqMTdybzB1b2MxdGgwbTYuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI1NDQzMDQxODc1OTYtMGxoYTRmZWwzZjY4N2szdTdqMTdybzB1b2MxdGgwbTYuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDIxMjM5NzI1MzA5MTQ0NzExODEiLCJlbWFpbCI6InZpY2JydW5vY3RvckBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6ImpUUzNNd2hNa3E5eHlCRmRGSk9nWlEiLCJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJpYXQiOjE1MTQ2NjAxMzksImV4cCI6MTUxNDY2MzczOSwibmFtZSI6IlZpY3RvciBCcnVubyIsInBpY3R1cmUiOiJodHRwczovL2xoNi5nb29nbGV1c2VyY29udGVudC5jb20vLVpwTUh5TWNaRk1nL0FBQUFBQUFBQUFJL0FBQUFBQUFBQUFBL0FGaVlvZjNLeDhMbzBXMkFwV3RsVGJmNlprSEtfSGVQOUEvczk2LWMvcGhvdG8uanBnIiwiZ2l2ZW5fbmFtZSI6IlZpY3RvciIsImZhbWlseV9uYW1lIjoiQnJ1bm8iLCJsb2NhbGUiOiJmciJ9.s4mO06c3LJV1BTEcWMx9vFOAGPkgeItE-p4_IAjKJZ235ab-RFH9_iWEQ0EA2n3AMK66UNtpUuDUdzZ-f9yEn1P2_OKTSaFfsa75cdm5hg6BfmxbZDLlb6B7vLrniLEun2kMMPTEXiV0EmfR7JOSiiuqWxtFdwP4XIvyJMqzGUHbE59tcH2VWSuTScAexuqdq78FedfJ-rpC_gswOs-o2r8bjiuPn8NnEfqWPNrUZnWArJ_cqYa5oHzh9WyX-QEy_tdYk-ziMQPIpple8ElCFUu_WktKdf-w_3DZg4g3ZqRT4f_wrOv6u8SalMM3w85nL3p3-xpzTYI0lYHBt7-K9A"
	Audience   = "544304187596-0lha4fel3f687k3u7j17ro0uoc1th0m6.apps.googleusercontent.com"
)

func TestSharedVerifier(t *testing.T) {
	claims, err := Decode(ValidToken, Audience)
	if err != nil {
		t.Error(err)
		return
	}

	if claims == nil {
		t.Errorf("Claims are nil")
	}
}

func TestCustomVerifier(t *testing.T) {
	v := NewVerifier().LazyLoading(false)
	if v.LazyLoad {
		t.Errorf("Lazy loading enabled, expected to be disabled")
	}

	v.LazyLoading(true)
	if !v.LazyLoad {
		t.Errorf("Lazy loading disabled, expected to be enabled")
	}

	v.PeriodicRefresh(time.Second * 1)
	if v.LazyLoad {
		t.Errorf("Lazy loading enabled, expected to be disabled by PeriodicRefresh")
	}

	v.LazyLoading(false)
}
