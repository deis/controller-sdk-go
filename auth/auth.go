// Package auth handles user management: creation, deletion, and authentication.
package auth

import (
	"encoding/json"
	"io"
	"io/ioutil"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// Register a new user with the controller.
// If controller registration is set to administratiors only, a valid administrative
// user token is required in the client.
func Register(c *deis.Client, username, password, email string) error {
	user := api.AuthRegisterRequest{Username: username, Password: password, Email: email}
	body, err := json.Marshal(user)

	if err != nil {
		return err
	}

	res, err := c.Request("POST", "/v2/auth/register/", body)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Login to the controller and get a token
func Login(c *deis.Client, username, password string) (string, error) {
	user := api.AuthLoginRequest{Username: username, Password: password}
	reqBody, err := json.Marshal(user)

	if err != nil {
		return "", err
	}

	res, err := c.Request("POST", "/v2/auth/login/", reqBody)
	if err != nil {
		return "", err
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	token := api.AuthLoginResponse{}
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return "", err
	}

	return token.Token, nil
}

// Delete deletes a user.
func Delete(c *deis.Client, username string) error {
	var body []byte
	var err error

	if username != "" {
		req := api.AuthCancelRequest{Username: username}
		body, err = json.Marshal(req)

		if err != nil {
			return err
		}
	}

	res, err := c.Request("DELETE", "/v2/auth/cancel/", body)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Regenerate auth tokens. This invalidates existing tokens, and if targeting a specific user
// returns a new token.
//
// If username is an empty string and all is false, this regenerates the
// client user's token and will return a new token. Make sure to update the client token
// with this new token to avoid authentication errors.
//
// If username is set and all is false, this will regenerate that user's token
// and return a new token. If not targeting yourself, regenerate requires administrative privilages.
//
// If all is true, this will regenerate every user's token. This requires administrative privilages.
func Regenerate(c *deis.Client, username string, all bool) (string, error) {
	var reqBody []byte
	var err error

	if all == true {
		reqBody, err = json.Marshal(api.AuthRegenerateRequest{All: all})
	} else if username != "" {
		reqBody, err = json.Marshal(api.AuthRegenerateRequest{Name: username})
	}

	if err != nil {
		return "", err
	}

	res, err := c.Request("POST", "/v2/auth/tokens/", reqBody)
	if err != nil {
		return "", err
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	if all == true {
		return "", nil
	}

	token := api.AuthRegenerateResponse{}
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return "", err
	}

	return token.Token, nil
}

// Passwd changes a user's password.
//
// If username if an empty string, change the password of the client's user.
//
// If username is set, change the password of another user and do not require
// their password. This requires administrative privilages.
func Passwd(c *deis.Client, username, password, newPassword string) error {
	req := api.AuthPasswdRequest{Password: password, NewPassword: newPassword}

	if username != "" {
		req.Username = username
	}

	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, err := c.Request("POST", "/v2/auth/passwd/", body)
	if err == nil {
		res.Body.Close()
	}
	return err
}

func Whoami(c *deis.Client) (api.User, error) {
	res, err := c.Request("GET", "/v2/auth/whoami/", nil)
	if err != nil {
		return api.User{}, err
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	resUser := api.User{}
	if err = json.NewDecoder(res.Body).Decode(&resUser); err != nil {
		return api.User{}, err
	}

	return resUser, nil
}
