// Package releases provides methods for managing app releases.
package releases

import (
	"encoding/json"
	"fmt"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List lists an app's releases.
func List(c *deis.Client, appID string, results int) ([]api.Release, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/releases/", appID)

	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return []api.Release{}, -1, reqErr
	}

	var releases []api.Release
	if err := json.Unmarshal([]byte(body), &releases); err != nil {
		return []api.Release{}, -1, err
	}

	return releases, count, reqErr
}

// Get retrieves a release of an app.
func Get(c *deis.Client, appID string, version int) (api.Release, error) {
	u := fmt.Sprintf("/v2/apps/%s/releases/v%d/", appID, version)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return api.Release{}, reqErr
	}
	defer res.Body.Close()

	release := api.Release{}
	if err := json.NewDecoder(res.Body).Decode(&release); err != nil {
		return api.Release{}, err
	}

	return release, reqErr
}

// Rollback rolls back an app to a previous release. If version is -1, this rolls back to
// the previous release. Otherwise, roll back to the specified version.
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

	res, reqErr := c.Request("POST", u, reqBody)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return -1, reqErr
	}
	defer res.Body.Close()

	response := api.ReleaseRollback{}

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return -1, err
	}

	return response.Version, reqErr
}
