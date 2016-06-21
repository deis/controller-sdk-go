package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List lists an app's config.
func List(c *deis.Client, app string) (api.Config, error) {
	u := fmt.Sprintf("/v2/apps/%s/config/", app)

	res, err := c.Request("GET", u, nil)
	if err != nil {
		return api.Config{}, err
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	config := api.Config{}
	if err = json.NewDecoder(res.Body).Decode(&config); err != nil {
		return api.Config{}, err
	}

	return config, nil
}

// Set sets an app's config variables.
func Set(c *deis.Client, app string, config api.Config) (api.Config, error) {
	body, err := json.Marshal(config)

	if err != nil {
		return api.Config{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/config/", app)

	res, err := c.Request("POST", u, body)
	if err != nil {
		return api.Config{}, err
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	newConfig := api.Config{}
	if err = json.NewDecoder(res.Body).Decode(&newConfig); err != nil {
		return api.Config{}, err
	}

	return newConfig, nil
}
