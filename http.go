package deis

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// createHTTPClient creates a HTTP Client with proper SSL options.
func createHTTPClient(sslVerify bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: !sslVerify},
		DisableKeepAlives: true,
		Proxy:             http.ProxyFromEnvironment,
	}
	return &http.Client{Transport: tr}
}

// Request makes a HTTP request with the given method, relative URL, and body on the controller.
// It also sets the Authorization and Content-Type headers to properly authenticate and communicate
// API. This is primarily intended to use be used by the SDK itself, but could potentially be used elsewhere.
func (c *Client) Request(method string, path string, body []byte) (*http.Response, error) {
	url := *c.ControllerURL

	if strings.Contains(path, "?") {
		parts := strings.Split(path, "?")
		url.Path = parts[0]
		url.RawQuery = parts[1]
	} else {
		url.Path = path
	}

	req, err := http.NewRequest(method, url.String(), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if c.Token != "" {
		req.Header.Add("Authorization", "token "+c.Token)
	}

	if c.HooksToken != "" {
		req.Header.Add("X-Deis-Builder-Auth", c.HooksToken)
	}

	addUserAgent(&req.Header, c.UserAgent)

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	if err = checkForErrors(res); err != nil {
		return nil, err
	}

	apiVersion := res.Header.Get("DEIS_API_VERSION")

	// Update controller api and platform version
	c.ControllerAPIVersion = apiVersion
	setControllerVersion(c, res.Header)

	// Return results along with api compatibility error
	return res, checkAPICompatibility(apiVersion, APIVersion)
}

// LimitedRequest allows limiting the number of responses in a request.
func (c *Client) LimitedRequest(path string, results int) (string, int, error) {
	res, reqErr := c.Request("GET", path+"?limit="+strconv.Itoa(results), nil)

	if reqErr != nil && !IsErrAPIMismatch(reqErr) {
		return "", -1, reqErr
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", -1, err
	}

	r := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return "", -1, err
	}

	out, err := json.Marshal(r["results"].([]interface{}))

	if err != nil {
		return "", -1, err
	}

	return string(out), int(r["count"].(float64)), reqErr
}

// CheckConnection checks that the user is connected to a network and the URL points to a valid controller.
func (c *Client) CheckConnection() error {
	errorMessage := `%s does not appear to be a valid Deis controller.
Make sure that the Controller URI is correct, the server is running and
your deis version is correct.`

	// Make a request to /v2/ and expect a 401 respone
	req, err := http.NewRequest("GET", c.ControllerURL.String()+"/v2/", bytes.NewBuffer(nil))
	addUserAgent(&req.Header, c.UserAgent)

	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 401 {
		return fmt.Errorf(errorMessage, c.ControllerURL.String())
	}

	// Update controller api version
	apiVersion := res.Header.Get("DEIS_API_VERSION")
	c.ControllerAPIVersion = apiVersion
	setControllerVersion(c, res.Header)

	return checkAPICompatibility(apiVersion, APIVersion)
}

// Healthcheck can be called to see if the controller is healthy
func (c *Client) Healthcheck() error {
	// Make a request to /healthz and expect an ok HTTP response
	req, err := http.NewRequest("GET", c.ControllerURL.String()+"/healthz", bytes.NewBuffer(nil))
	addUserAgent(&req.Header, c.UserAgent)

	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return err
	}

	if err = checkForErrors(res); err != nil {
		return err
	}
	res.Body.Close()

	// Update controller api version
	apiVersion := res.Header.Get("DEIS_API_VERSION")
	c.ControllerAPIVersion = apiVersion
	setControllerVersion(c, res.Header)

	return checkAPICompatibility(apiVersion, APIVersion)
}

func addUserAgent(headers *http.Header, userAgent string) {
	headers.Add("User-Agent", userAgent)
}

func setControllerVersion(c *Client, headers http.Header) {
	c.ControllerVersion = headers.Get("DEIS_PLATFORM_VERSION")
}
