// Package whitelist provides methods for managing an app's whitelisted IP's.
package whitelist

import (
	"encoding/json"
	"fmt"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List IP's whitelisted for an app.
func List(c *deis.Client, appID string) (api.Whitelist, error) {
	u := fmt.Sprintf("/v2/apps/%s/whitelist/", appID)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return api.Whitelist{}, reqErr
	}
	defer res.Body.Close()

	whitelist := api.Whitelist{}
	if err := json.NewDecoder(res.Body).Decode(&whitelist); err != nil {
		return api.Whitelist{}, err
	}

	return whitelist, reqErr
}

// Add adds addresses to an app's whitelist.
func Add(c *deis.Client, appID string, addresses []string) (api.Whitelist, error) {
	u := fmt.Sprintf("/v2/apps/%s/whitelist/", appID)

	req := api.Whitelist{Addresses: addresses}
	body, err := json.Marshal(req)
	if err != nil {
		return api.Whitelist{}, err
	}
	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return api.Whitelist{}, reqErr
	}
	defer res.Body.Close()

	d := api.Whitelist{}
	if err = json.NewDecoder(res.Body).Decode(&d); err != nil {
		return api.Whitelist{}, err
	}

	return d, reqErr
}

// Delete removes addresses from an app's whitelist.
func Delete(c *deis.Client, appID string, addresses []string) error {
	u := fmt.Sprintf("/v2/apps/%s/whitelist/", appID)

	req := api.Whitelist{Addresses: addresses}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, reqErr := c.Request("DELETE", u, body)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return reqErr
	}
	return nil
}
