// Package tls provides methods for managing tls configuration for apps.
package tls

import (
	"encoding/json"
	"fmt"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// Info displays an app's tls config.
func Info(c *deis.Client, app string) (api.TLS, error) {
	u := fmt.Sprintf("/v2/apps/%s/tls/", app)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil {
		return api.TLS{}, reqErr
	}
	defer res.Body.Close()

	tls := api.TLS{}
	if err := json.NewDecoder(res.Body).Decode(&tls); err != nil {
		return api.TLS{}, err
	}

	return tls, reqErr
}

// Enable enables the router to enforce https-only requests to the application.
func Enable(c *deis.Client, app string) (api.TLS, error) {
	t := api.NewTLS()
	b := true
	t.HTTPSEnforced = &b
	body, err := json.Marshal(t)

	if err != nil {
		return api.TLS{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/tls/", app)

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil {
		return api.TLS{}, reqErr
	}
	defer res.Body.Close()

	newTLS := api.TLS{}
	if err = json.NewDecoder(res.Body).Decode(&newTLS); err != nil {
		return api.TLS{}, err
	}

	return newTLS, reqErr
}

// Disable disables the router from enforcing https-only requests to the application.
func Disable(c *deis.Client, app string) (api.TLS, error) {
	body, err := json.Marshal(api.NewTLS())

	if err != nil {
		return api.TLS{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/tls/", app)

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil {
		return api.TLS{}, reqErr
	}
	defer res.Body.Close()

	newTLS := api.TLS{}
	if err = json.NewDecoder(res.Body).Decode(&newTLS); err != nil {
		return api.TLS{}, err
	}

	return newTLS, reqErr
}
