package perms

import (
	"encoding/json"
	"fmt"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List users that can access an app.
func List(c *deis.Client, appID string) ([]string, error) {
	body, err := c.BasicRequest("GET", fmt.Sprintf("/v2/apps/%s/perms/", appID), nil)

	if err != nil {
		return []string{}, err
	}

	var users api.PermsAppResponse
	if err = json.Unmarshal([]byte(body), &users); err != nil {
		return []string{}, err
	}

	return users.Users, nil
}

// ListAdmins lists administrators.
func ListAdmins(c *deis.Client, results int) ([]string, int, error) {
	body, count, err := c.LimitedRequest("/v2/admin/perms/", results)

	if err != nil {
		return []string{}, -1, err
	}

	var users []api.PermsRequest
	if err = json.Unmarshal([]byte(body), &users); err != nil {
		return []string{}, -1, err
	}

	usersList := []string{}

	for _, user := range users {
		usersList = append(usersList, user.Username)
	}

	return usersList, count, nil
}

// New adds a user to an app.
func New(c *deis.Client, appID string, username string) error {
	return doNew(c, fmt.Sprintf("/v2/apps/%s/perms/", appID), username)
}

// NewAdmin makes a user an administrator.
func NewAdmin(c *deis.Client, username string) error {
	return doNew(c, "/v2/admin/perms/", username)
}

func doNew(c *deis.Client, u string, username string) error {
	req := api.PermsRequest{Username: username}

	reqBody, err := json.Marshal(req)

	if err != nil {
		return err
	}

	_, err = c.BasicRequest("POST", u, reqBody)

	if err != nil {
		return err
	}

	return nil
}

// Delete removes a user from an app.
func Delete(c *deis.Client, appID string, username string) error {
	return doDelete(c, fmt.Sprintf("/v2/apps/%s/perms/%s", appID, username))
}

// DeleteAdmin removes administrative privileges from a user.
func DeleteAdmin(c *deis.Client, username string) error {
	return doDelete(c, fmt.Sprintf("/v2/admin/perms/%s", username))
}

func doDelete(c *deis.Client, u string) error {
	_, err := c.BasicRequest("DELETE", u, nil)
	return err
}
