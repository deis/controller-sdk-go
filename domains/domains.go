// Package domains provides methods for managing an app's domains.
package domains

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List domains registered with an app.
func List(c *deis.Client, appID string, results int) (api.Domains, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/domains/", appID)
	body, count, err := c.LimitedRequest(u, results)

	if err != nil {
		return []api.Domain{}, -1, err
	}

	var domains []api.Domain
	if err = json.Unmarshal([]byte(body), &domains); err != nil {
		return []api.Domain{}, -1, err
	}

	return domains, count, nil
}

// New adds a domain to an app.
func New(c *deis.Client, appID string, domain string) (api.Domain, error) {
	u := fmt.Sprintf("/v2/apps/%s/domains/", appID)

	req := api.DomainCreateRequest{Domain: domain}

	body, err := json.Marshal(req)

	if err != nil {
		return api.Domain{}, err
	}

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return api.Domain{}, reqErr
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	d := api.Domain{}
	if err = json.NewDecoder(res.Body).Decode(&d); err != nil {
		return api.Domain{}, err
	}

	return d, reqErr
}

// Delete removes a domain from an app.
func Delete(c *deis.Client, appID string, domain string) error {
	u := fmt.Sprintf("/v2/apps/%s/domains/%s", appID, domain)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
