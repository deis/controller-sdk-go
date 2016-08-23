// Package appsettings provides methods for managing application settings of apps.
package appsettings

import (
	"encoding/json"
	"fmt"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List lists an app's settings.
func List(c *deis.Client, app string) (api.AppSettings, error) {
	u := fmt.Sprintf("/v2/apps/%s/settings/", app)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil {
		return api.AppSettings{}, reqErr
	}
	defer res.Body.Close()

	settings := api.AppSettings{}
	if err := json.NewDecoder(res.Body).Decode(&settings); err != nil {
		return api.AppSettings{}, err
	}

	return settings, reqErr
}

// Set sets an app's settings variables.
// This is a patching operation, which means when you call Set() with an api.AppSettings:
//
//    - If the variable does not exist, it will be set.
//    - If the variable exists, it will be overwritten.
//    - If the variable is set to nil, it will be unset.
//    - If the variable was ignored in the api.AppSettings, it will remain unchanged.
//
// Calling Set() with an empty api.AppSettings will return a deis.ErrConflict.
func Set(c *deis.Client, app string, appSettings api.AppSettings) (api.AppSettings, error) {
	body, err := json.Marshal(appSettings)

	if err != nil {
		return api.AppSettings{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/settings/", app)

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil {
		return api.AppSettings{}, reqErr
	}
	defer res.Body.Close()

	newAppSettings := api.AppSettings{}
	if err = json.NewDecoder(res.Body).Decode(&newAppSettings); err != nil {
		return api.AppSettings{}, err
	}

	return newAppSettings, reqErr
}
