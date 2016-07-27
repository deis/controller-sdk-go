package users

import (
	"encoding/json"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List users registered with the controller.
func List(c *deis.Client, results int) (api.Users, int, error) {
	body, count, err := c.LimitedRequest("/v2/users/", results)

	if err != nil {
		return []api.User{}, -1, err
	}

	var users []api.User
	if err = json.Unmarshal([]byte(body), &users); err != nil {
		return []api.User{}, -1, err
	}

	return users, count, nil
}
