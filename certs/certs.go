package certs

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List certs registered with the controller.
func List(c *deis.Client, results int) ([]api.Cert, int, error) {
	body, count, err := c.LimitedRequest("/v2/certs/", results)

	if err != nil {
		return []api.Cert{}, -1, err
	}

	var res []api.Cert
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return []api.Cert{}, -1, err
	}

	return res, count, nil
}

// New creates a cert.
func New(c *deis.Client, cert string, key string, name string) (api.Cert, error) {
	req := api.CertCreateRequest{Certificate: cert, Key: key, Name: name}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return api.Cert{}, err
	}

	res, reqErr := c.Request("POST", "/v2/certs/", reqBody)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return api.Cert{}, reqErr
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	resCert := api.Cert{}
	if err = json.NewDecoder(res.Body).Decode(&resCert); err != nil {
		return api.Cert{}, err
	}

	return resCert, reqErr
}

// Get information for a certificate
func Get(c *deis.Client, name string) (api.Cert, error) {
	url := fmt.Sprintf("/v2/certs/%s", name)
	res, reqErr := c.Request("GET", url, nil)
	if reqErr != nil {
		return api.Cert{}, reqErr
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	resCert := api.Cert{}
	if err := json.NewDecoder(res.Body).Decode(&resCert); err != nil {
		return api.Cert{}, err
	}

	return resCert, reqErr
}

// Delete removes a cert.
func Delete(c *deis.Client, name string) error {
	url := fmt.Sprintf("/v2/certs/%s", name)
	res, err := c.Request("DELETE", url, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Attach a certificate to a domain
func Attach(c *deis.Client, name string, domain string) error {
	req := api.CertAttachRequest{Domain: domain}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/v2/certs/%s/domain/", name)
	res, err := c.Request("POST", url, reqBody)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Detach a certificate from a domain
func Detach(c *deis.Client, name string, domain string) error {
	url := fmt.Sprintf("/v2/certs/%s/domain/%s", name, domain)
	res, err := c.Request("DELETE", url, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
