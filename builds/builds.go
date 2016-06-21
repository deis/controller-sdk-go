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
	body, count, err := c.LimitedRequest(u, results)

	if err != nil {
		return []api.Build{}, -1, err
	}

	var builds []api.Build
	if err = json.Unmarshal([]byte(body), &builds); err != nil {
		return []api.Build{}, -1, err
	}

	return builds, count, nil
}

// New creates a build for an app.
func New(c *deis.Client, appID string, image string,
	procfile map[string]string) (api.Build, error) {

	u := fmt.Sprintf("/v2/apps/%s/builds/", appID)

	req := api.CreateBuildRequest{Image: image, Procfile: procfile}

	body, err := json.Marshal(req)

	if err != nil {
		return api.Build{}, err
	}

	res, err := c.Request("POST", u, body)
	if err != nil {
		return api.Build{}, err
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

	return build, nil
}
