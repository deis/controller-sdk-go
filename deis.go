// Package deis offers a SDK for interacting with the Deis controller API.
//
// This package works by creating a client, which contains session information,
// such as the controller url and user token. The client is then passed to api methods,
// which use it to make requests.
//
// Basic Example
//
// This example creates a client and then lists the apps that the user has access to:
//
//    import (
//        deis "github.com/deis/controller-sdk-go"
//        "github.com/deis/controller-sdk-go/apps"
//    )
//
//    //                      Verify SSL, Controller URL, API Token
//    client, err := deis.New(true, "deis.test.io", "abc123")
//    if err != nil {
//        log.Fatal(err)
//    }
//    apps, _, err := apps.List(client, 100)
//    if err != nil {
//        log.Fatal(err)
//    }
//
// Authentication
//
// If you don't already have a token for a user, you can retrieve one with a
// username and password.
//
//    import (
//        deis "github.com/deis/controller-sdk-go"
//        "github.com/deis/controller-sdk-go/apps"
//    )
//
//    // Create a client with a blank token to pass to login.
//    client, err := deis.New(true, "deis.test.io", "")
//    if err != nil {
//        log.Fatal(err)
//    }
//    token, err := auth.Login(client, "user", "password")
//    if err != nil {
//        log.Fatal(err)
//    }
//    // Set the client to use the retrieved token
//    client.Token = token
//
// Learning More
//
// See the godoc for the SDK's subpackages to learn more about specific SDK actions.
package deis

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/goware/urlx"
)

// Client oversees the interaction between the deis and controller
type Client struct {
	// HTTPClient is the transport that is used to communicate with the API.
	HTTPClient *http.Client

	// VerifySSL determines whether or not to verify SSL connections.
	// This should be true unless you know the controller is using untrusted SSL keys.
	VerifySSL bool

	// ControllerURL is the URL used to communicate with the controller.
	ControllerURL *url.URL

	// UserAgent is the user agent used when making requests.
	UserAgent string

	// API Version used by the controller, set after a http request.
	ControllerAPIVersion string

	// Version of the deis platform in use, set after a http request.
	DeisVersion string

	// Token is used to authenticate the request against the API.
	Token string
}

// APIVersion is the api version compatible with the SDK.
//
// In general, using an SDK that is a minor version out of date with the target controller
// is probably safe, as the deis controller api follows semantic versioning and is backward
// compatible. However, using a SDK that is newer or a major version different than the
// controller is unsafe.
//
// If the SDK detects an API version mismatch, it will return ErrAPIMismatch.
const APIVersion = "2.2"

var (
	// ErrAPIMismatch occurs when the sdk is using a different api version than the deis.
	ErrAPIMismatch = errors.New("API Version Mismatch between server and deis")

	// DefaultUserAgent is used as the default user agent when making requests.
	DefaultUserAgent = fmt.Sprintf("Deis Go SDK V%s", APIVersion)
)

// IsErrAPIMismatch returns true if err is an ErrAPIMismatch, false otherwise
func IsErrAPIMismatch(err error) bool {
	return err == ErrAPIMismatch
}

// New creates a new client to communicate with the api.
// The controllerURL is the url of the controller component, by default deis.<cluster url>.com
// verifySSL determines whether or not to verify SSL connections.
// This should be true unless you know the controller is using untrusted SSL keys.
func New(verifySSL bool, controllerURL string, token string) (*Client, error) {
	// urlx, unlike the native url library, uses sane defaults when URL parsing,
	// preventing issues like missing schemes.
	u, err := urlx.Parse(controllerURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		HTTPClient:    createHTTPClient(verifySSL),
		VerifySSL:     verifySSL,
		ControllerURL: u,
		Token:         token,
		UserAgent:     DefaultUserAgent,
	}, nil
}
