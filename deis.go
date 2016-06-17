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

	// Username is the name of the user performing requests against the API.
	Username string

	// ResponseLimit is the number of results to return on requests that can be limited.
	ResponseLimit int
}

// APIVersion is the api version the sdk is compatible with.
const APIVersion = "2.0"

var (
	// ErrAPIMismatch occurs when the sdk is using a different api version than the deis.
	ErrAPIMismatch = errors.New("API Version Mismatch between server and deis")

	// DefaultResponseLimit is the default number of responses to return on requests that can
	// be limited.
	DefaultResponseLimit = 100

	// DefaultUserAgent is used as the default user agent when making requests.
	DefaultUserAgent = fmt.Sprintf("Deis Go SDK V%s", APIVersion)
)

// New creates a new deis to communicate with the api.
func New(verifySSL bool, controllerURL string, token string, username string) (*Client, error) {
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
		Username:      username,
		ResponseLimit: DefaultResponseLimit,
		UserAgent:     DefaultUserAgent,
	}, nil
}
