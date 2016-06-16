package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
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

// Request makes a HTTP request on the controller.
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

	addUserAgent(&req.Header, c.UserAgent)

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	if err = checkForErrors(res, ""); err != nil {
		return nil, err
	}

	apiVersion := res.Header.Get("DEIS_API_VERSION")

	// Update controller api version
	c.ControllerAPIVersion = apiVersion

	// Return results along with api compatibility error
	return res, checkAPICompatibility(apiVersion)
}

// LimitedRequest allows limiting the number of responses in a request.
func (c *Client) LimitedRequest(path string, results int) (string, int, error) {
	body, err := c.BasicRequest("GET", path+"?limit="+strconv.Itoa(results), nil)

	if err != nil {
		return "", -1, err
	}

	res := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return "", -1, err
	}

	out, err := json.Marshal(res["results"].([]interface{}))

	if err != nil {
		return "", -1, err
	}

	return string(out), int(res["count"].(float64)), nil
}

// BasicRequest makes a simple http request on the controller.
func (c *Client) BasicRequest(method string, path string, body []byte) (string, error) {
	res, err := c.Request(method, path, body)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}
	return string(resBody), checkForErrors(res, string(resBody))
}

func checkForErrors(res *http.Response, body string) error {

	// If response is not an error, return nil.
	if res.StatusCode > 199 && res.StatusCode < 400 {
		return nil
	}

	// Read the response body if none was provided.
	if body == "" {
		defer res.Body.Close()
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		body = string(resBody)
	}

	// Unmarshal the response as JSON, or return the status and body.
	bodyMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &bodyMap); err != nil {
		return fmt.Errorf("\n%s\n%s\n", res.Status, body)
	}

	errorMessage := fmt.Sprintf("\n%s\n", res.Status)
	for key, value := range bodyMap {
		switch v := value.(type) {
		case string:
			errorMessage += fmt.Sprintf("%s: %s\n", key, v)
		case []interface{}:
			for _, subValue := range v {
				switch sv := subValue.(type) {
				case string:
					errorMessage += fmt.Sprintf("%s: %s\n", key, sv)
				default:
					fmt.Printf("Unexpected type in %s error message array. Contents: %v",
						reflect.TypeOf(value), sv)
				}
			}
		default:
			fmt.Printf("Cannot handle key %s in error message, type %s. Contents: %v",
				key, reflect.TypeOf(value), bodyMap[key])
		}
	}

	return errors.New(errorMessage)
}

// CheckConnection checks that the user is connected to a network and the URL points to a valid controller.
func (c *Client) CheckConnection() error {
	errorMessage := `%s does not appear to be a valid Deis controller.
Make sure that the Controller URI is correct, the server is running and
your client version is correct.`

	// Make a request to /v2/ and expect a 401 respone
	req, err := http.NewRequest("GET", c.ControllerURL.String()+"/v2/", bytes.NewBuffer(nil))
	addUserAgent(&req.Header, c.UserAgent)

	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		fmt.Printf(errorMessage+"\n", c.ControllerURL.String())
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 401 {
		return fmt.Errorf(errorMessage, c.ControllerURL.String())
	}

	// Update controller api version
	apiVersion := res.Header.Get("DEIS_API_VERSION")
	c.ControllerAPIVersion = apiVersion

	return checkAPICompatibility(apiVersion)
}

func addUserAgent(headers *http.Header, userAgent string) {
	headers.Add("User-Agent", userAgent)
}
