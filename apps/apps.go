// Package apps provides methods for managing deis apps.
package apps

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

const workflowURLPrefix = "deis."

// ErrNoLogs is returned when logs are missing from an app.
var ErrNoLogs = errors.New(
	`There are currently no log messages. Please check the following things:
1) Logger and fluentd pods are running: kubectl --namespace=deis get pods.
2) The application is writing logs to the logger component by checking that an entry in the ring buffer was created: kubectl --namespace=deis logs <logger pod>
3) Making sure that the container logs were mounted properly into the fluentd pod: kubectl --namespace=deis exec <fluentd pod> ls /var/log/containers
3a) If the above command returns saying /var/log/containers cannot be found then please see the following github issue for a workaround: https://github.com/deis/logger/issues/50`)

// List lists apps on a Deis controller.
func List(c *deis.Client, results int) (api.Apps, int, error) {
	body, count, err := c.LimitedRequest("/v2/apps/", results)

	if err != nil {
		return []api.App{}, -1, err
	}

	var apps []api.App
	if err = json.Unmarshal([]byte(body), &apps); err != nil {
		return []api.App{}, -1, err
	}

	for name, app := range apps {
		// Add in app URL based on controller hostname, port included
		app.URL = fmt.Sprintf("%s.%s", app.ID, strings.TrimPrefix(c.ControllerURL.Host, workflowURLPrefix))
		apps[name] = app
	}

	return apps, count, nil
}

// New creates a new app with the given appID. Passing an empty string will result in
// a randomized app name.
//
// If the app name already exists, the error deis.ErrDuplicateApp will be returned.
func New(c *deis.Client, appID string) (api.App, error) {
	body := []byte{}

	if appID != "" {
		req := api.AppCreateRequest{ID: appID}
		b, err := json.Marshal(req)

		if err != nil {
			return api.App{}, err
		}
		body = b
	}

	res, reqErr := c.Request("POST", "/v2/apps/", body)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return api.App{}, reqErr
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	app := api.App{}
	if err := json.NewDecoder(res.Body).Decode(&app); err != nil {
		return api.App{}, err
	}

	// Add in app URL based on controller hostname, port included
	app.URL = fmt.Sprintf("%s.%s", app.ID, strings.TrimPrefix(c.ControllerURL.Host, workflowURLPrefix))

	return app, reqErr
}

// Get app details from a controller.
func Get(c *deis.Client, appID string) (api.App, error) {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return api.App{}, reqErr
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	app := api.App{}

	if err := json.NewDecoder(res.Body).Decode(&app); err != nil {
		return api.App{}, err
	}

	// Add in app URL based on controller hostname, port included
	app.URL = fmt.Sprintf("%s.%s", app.ID, strings.TrimPrefix(c.ControllerURL.Host, workflowURLPrefix))

	return app, reqErr
}

// Logs retrieves logs from an app. The number of log lines fetched can be set by the lines
// argument. Setting lines = -1 will retrive all app logs.
func Logs(c *deis.Client, appID string, lines int) (string, error) {
	u := fmt.Sprintf("/v2/apps/%s/logs", appID)

	if lines > 0 {
		u += "?log_lines=" + strconv.Itoa(lines)
	}

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return "", ErrNoLogs
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil || len(body) < 3 {
		return "", ErrNoLogs
	}

	// We need to trim a few characters off the front and end of the string
	return string(body[2 : len(body)-1]), reqErr
}

// Run a one-time command in your app. This will start a kubernetes job with the
// same container image and environment as the rest of the app.
func Run(c *deis.Client, appID string, command string) (api.AppRunResponse, error) {
	req := api.AppRunRequest{Command: command}
	body, err := json.Marshal(req)

	if err != nil {
		return api.AppRunResponse{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/run", appID)

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return api.AppRunResponse{}, reqErr
	}

	arr := api.AppRunResponse{}

	if err = json.NewDecoder(res.Body).Decode(&arr); err != nil {
		return api.AppRunResponse{}, err
	}

	return arr, reqErr
}

// Delete an app.
func Delete(c *deis.Client, appID string) error {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Transfer an app to another user.
func Transfer(c *deis.Client, appID string, username string) error {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	req := api.AppUpdateRequest{Owner: username}
	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, err := c.Request("POST", u, body)
	if err == nil {
		res.Body.Close()
	}
	return err
}
