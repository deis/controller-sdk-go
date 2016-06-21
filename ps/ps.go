package ps

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List an app's processes.
func List(c *deis.Client, appID string, results int) ([]api.Pods, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/pods/", appID)
	body, count, err := c.LimitedRequest(u, results)
	if err != nil {
		return []api.Pods{}, -1, err
	}

	var procs []api.Pods
	if err = json.Unmarshal([]byte(body), &procs); err != nil {
		return []api.Pods{}, -1, err
	}

	return procs, count, nil
}

// Scale an app's processes.
func Scale(c *deis.Client, appID string, targets map[string]int) error {
	u := fmt.Sprintf("/v2/apps/%s/scale/", appID)

	body, err := json.Marshal(targets)

	if err != nil {
		return err
	}

	res, err := c.Request("POST", u, body)
	if err == nil {
		return res.Body.Close()
	}
	return err
}

// Restart an app's processes.
func Restart(c *deis.Client, appID string, procType string, name string) ([]api.Pods, error) {
	u := fmt.Sprintf("/v2/apps/%s/pods/", appID)

	if procType == "" {
		u += "restart/"
	} else {
		if name == "" {
			u += procType + "/restart/"
		} else {
			u += procType + "/" + name + "/restart/"
		}
	}

	res, err := c.Request("POST", u, nil)
	if err != nil {
		return []api.Pods{}, err
	}
	// Fix json.Decoder bug in <go1.7
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	procs := []api.Pods{}
	if err = json.NewDecoder(res.Body).Decode(&procs); err != nil {
		return []api.Pods{}, err
	}

	return procs, nil
}

// ByType organizes processes of an app by process type.
func ByType(processes []api.Pods) map[string][]api.Pods {
	psMap := make(map[string][]api.Pods)

	for _, ps := range processes {
		psMap[ps.Type] = append(psMap[ps.Type], ps)
	}

	return psMap
}
