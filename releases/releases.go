package releases

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List lists an app's releases.
func List(c *deis.Client, appID string, results int) ([]api.Release, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/releases/", appID)

	body, count, err := c.LimitedRequest(u, results)

	if err != nil {
		return []api.Release{}, -1, err
	}

	var releases []api.Release
	if err = json.Unmarshal([]byte(body), &releases); err != nil {
		return []api.Release{}, -1, err
	}

	return releases, count, nil
}

// Get a release of an app.
func Get(c *deis.Client, appID string, version int) (api.Release, error) {
	u := fmt.Sprintf("/v2/apps/%s/releases/v%d/", appID, version)

	res, err := c.Request("GET", u, nil)
	if err != nil {
		return api.Release{}, err
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	release := api.Release{}
	if err = json.NewDecoder(res.Body).Decode(&release); err != nil {
		return api.Release{}, err
	}

	return release, nil
}

// Rollback rolls back an app to a previous release.
func Rollback(c *deis.Client, appID string, version int) (int, error) {
	u := fmt.Sprintf("/v2/apps/%s/releases/rollback/", appID)

	req := api.ReleaseRollback{Version: version}

	var err error
	var reqBody []byte
	if version != -1 {
		reqBody, err = json.Marshal(req)

		if err != nil {
			return -1, err
		}
	}

	res, err := c.Request("POST", u, reqBody)
	if err != nil {
		return -1, err
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	response := api.ReleaseRollback{}

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return -1, err
	}

	return response.Version, nil
}
