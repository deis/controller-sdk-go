// Package users provides methods for viewing users.
package users

import (
	"encoding/json"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List lists users registered with the controller.
func List(c *deis.Client, results int) (api.Users, int, error) {
	body, count, reqErr := c.LimitedRequest("/v2/users/", results)

	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return []api.User{}, -1, reqErr
	}

	var users []api.User
	if err := json.Unmarshal([]byte(body), &users); err != nil {
		return []api.User{}, -1, err
	}

	return users, count, reqErr
}
