package api

import (
	"strings"
	"testing"
)

func TestTLSString(t *testing.T) {
	tls := &TLS{}

	expected := "HTTPS Enforced: not set"

	if strings.TrimSpace(tls.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, tls.String())
	}

	tls = NewTLS()

	expected = "HTTPS Enforced: false"

	if strings.TrimSpace(tls.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, tls.String())
	}

	b := true
	tls.HTTPSEnforced = &b

	expected = "HTTPS Enforced: true"

	if strings.TrimSpace(tls.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, tls.String())
	}
}
