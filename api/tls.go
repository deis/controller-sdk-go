package api

import (
	"fmt"
)

// TLS is the structure of an app's TLS settings.
type TLS struct {
	// Owner is the app owner. It cannot be updated with TLS.Set(). See app.Transfer().
	Owner string `json:"owner,omitempty"`
	// App is the app the tls settings apply to and cannot be updated.
	App string `json:"app,omitempty"`
	// Created is the time that the TLS settings was created and cannot be updated.
	Created string `json:"created,omitempty"`
	// Updated is the last time the TLS settings was changed and cannot be updated.
	Updated string `json:"updated,omitempty"`
	// UUID is a unique string reflecting the TLS settings in its current state.
	// It changes every time the TLS settings is changed and cannot be updated.
	UUID string `json:"uuid,omitempty"`
	//HTTPSEnforced determines if the router should enable or disable https-only requests.
	HTTPSEnforced *bool `json:"https_enforced,omitempty"`
}

// NewTLS creates a new TLS object with fields properly zeroed
func NewTLS() *TLS {
	return &TLS{
		HTTPSEnforced: new(bool),
	}
}

func (t TLS) String() string {
	if t.HTTPSEnforced != nil {
		return fmt.Sprintf("HTTPS Enforced: %t", *t.HTTPSEnforced)
	}
	return "HTTPS Enforced: not set"
}
