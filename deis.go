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
	// HTTP deis used to communicate with the API.
	HTTPClient *http.Client

	// VerifySSL determines whether or not to verify SSL connections.
	VerifySSL bool

	// URL used to communicate with the controller.
	ControllerURL *url.URL

	//UserAgent is the user agent used when making requests
	UserAgent string

	//API Version used by the controller, set after a http request.
	ControllerAPIVersion string

	// Token is used to authenticate the request against the API.
	Token string
}

// APIVersion is the api version the sdk is compatible with.
const APIVersion = "2.0"

var (
	// ErrAPIMismatch occurs when the sdk is using a different api version than the deis.
	ErrAPIMismatch = errors.New("API Version Mismatch between server and deis")

	// DefaultUserAgent is used as the default user agent when making requests.
	DefaultUserAgent = fmt.Sprintf("Deis Go SDK V%s", APIVersion)
)

// New creates a new deis to communicate with the api.
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
