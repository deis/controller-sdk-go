// Package config provides methods for managing configuration of apps.
package config

import (
	"encoding/json"
	"fmt"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List lists an app's config.
func List(c *deis.Client, app string) (api.Config, error) {
	u := fmt.Sprintf("/v2/apps/%s/config/", app)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil {
		return api.Config{}, reqErr
	}
	defer res.Body.Close()

	config := api.Config{}
	if err := json.NewDecoder(res.Body).Decode(&config); err != nil {
		return api.Config{}, err
	}

	return config, reqErr
}

// Set sets an app's config variables and creates a new release.
// This is a patching operation, which means when you call Set() with an api.Config:
//
//    - If the variable does not exist, it will be set.
//    - If the variable exists, it will be overwritten.
//    - If the variable is set to nil, it will be unset.
//    - If the variable was ignored in the api.Config, it will remain unchanged.
//
// Calling Set() with an empty api.Config will return a deis.ErrConflict.
// Trying to unset a key that does not exist returns a deis.ErrUnprocessable.
// Trying to set a tag that is not a label in the kubernetes cluster will return a deis.ErrTagNotFound.
func Set(c *deis.Client, app string, config api.Config) (api.Config, error) {
	body, err := json.Marshal(config)

	if err != nil {
		return api.Config{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/config/", app)

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil {
		return api.Config{}, reqErr
	}
	defer res.Body.Close()

	newConfig := api.Config{}
	if err = json.NewDecoder(res.Body).Decode(&newConfig); err != nil {
		return api.Config{}, err
	}

	return newConfig, reqErr
}
