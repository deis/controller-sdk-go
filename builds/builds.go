// Package builds provides methods for managing app builds.
package builds

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List lists an app's builds.
func List(c *deis.Client, appID string, results int) ([]api.Build, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/builds/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return []api.Build{}, -1, reqErr
	}

	var builds []api.Build
	if err := json.Unmarshal([]byte(body), &builds); err != nil {
		return []api.Build{}, -1, err
	}

	return builds, count, reqErr
}

// New creates a build for an app from an docker image.
// By default this will create a cmd process that runs the CMD command from the Dockerfile.
// If you want to define more process types, you can pass a Procfile map,
// where the key is the process name and the value is the command for that process.
// To pull from a private docker registry, a custom username and password must be set in the app's
// configuration object. This can be done with `deis registry:set` or by using this SDK.
//
// This example adds custom registry credentials to an app:
//    import (
//    	"github.com/deis/controller-sdk-go/api"
//    	"github.com/deis/controller-sdk-go/config"
//    )
//
//    // Create username/password map
//    registryMap := map[string]string{
//    	"username": "password"
//    }
//
//    // Create a new configuration, assign the credentials, and set it.
//    // Note that config setting is a patching operation, it doesn't overwrite or unset
//    // unrelated configuration.
//    newConfig := api.Config{}
//    newConfig.Registry = registryMap
//    _, err := config.Set(<client>, "appname", newConfig)
//    if err != nil {
//        log.Fatal(err)
//    }
func New(c *deis.Client, appID string, image string,
	procfile map[string]string) (api.Build, error) {

	u := fmt.Sprintf("/v2/apps/%s/builds/", appID)

	req := api.CreateBuildRequest{Image: image, Procfile: procfile}

	body, err := json.Marshal(req)

	if err != nil {
		return api.Build{}, err
	}

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return api.Build{}, reqErr
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	build := api.Build{}
	if err = json.NewDecoder(res.Body).Decode(&build); err != nil {
		return api.Build{}, err
	}

	return build, reqErr
}
