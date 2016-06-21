package keys

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List keys on a controller.
func List(c *deis.Client, results int) ([]api.Key, int, error) {
	body, count, err := c.LimitedRequest("/v2/keys/", results)

	if err != nil {
		return []api.Key{}, -1, err
	}

	var keys []api.Key
	if err = json.Unmarshal([]byte(body), &keys); err != nil {
		return []api.Key{}, -1, err
	}

	return keys, count, nil
}

// New creates a new key.
func New(c *deis.Client, id string, pubKey string) (api.Key, error) {
	req := api.KeyCreateRequest{ID: id, Public: pubKey}
	body, err := json.Marshal(req)

	res, err := c.Request("POST", "/v2/keys/", body)
	if err != nil {
		return api.Key{}, err
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	key := api.Key{}
	if err = json.NewDecoder(res.Body).Decode(&key); err != nil {
		return api.Key{}, err
	}

	return key, nil
}

// Delete a key.
func Delete(c *deis.Client, keyID string) error {
	u := fmt.Sprintf("/v2/keys/%s", keyID)

	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
