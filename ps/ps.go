// Package ps provides methods for managing app processes.
package ps

import (
	"encoding/json"
	"fmt"
	"sort"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

// List lists an app's processes.
func List(c *deis.Client, appID string, results int) (api.PodsList, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/pods/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return []api.Pods{}, -1, reqErr
	}

	var procs []api.Pods
	if err := json.Unmarshal([]byte(body), &procs); err != nil {
		return []api.Pods{}, -1, err
	}

	return procs, count, reqErr
}

// Scale increases or decreases an app's processes. The processes are specified in the target argument,
// a key-value map, where the key is the process name and the value is the number of replicas
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

// Restart restarts an app's processes. To restart all app processes, pass empty strings for
// procType and name. To restart an specific process, pass an procType by leave name empty.
// To restart a specific instance, pass a procType and a name.
func Restart(c *deis.Client, appID string, procType string, name string) (api.PodsList, error) {
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

	res, reqErr := c.Request("POST", u, nil)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return []api.Pods{}, reqErr
	}
	defer res.Body.Close()

	procs := []api.Pods{}
	if err := json.NewDecoder(res.Body).Decode(&procs); err != nil {
		return []api.Pods{}, err
	}

	return procs, reqErr
}

// ByType organizes processes of an app by process type.
func ByType(processes api.PodsList) api.PodTypes {
	var pts api.PodTypes

	for _, process := range processes {
		exists := false
		// Is processtype for process already exists, append to it.
		for i, pt := range pts {
			if pt.Type == process.Type {
				exists = true
				pts[i].PodsList = append(pts[i].PodsList, process)
				break
			}
		}

		// Is processtype for process doesn't exist, create a new one
		if !exists {
			pts = append(pts, api.PodType{
				Type:     process.Type,
				PodsList: api.PodsList{process},
			})
		}
	}

	// Sort the pods alphabetically by name.
	for _, pt := range pts {
		sort.Sort(pt.PodsList)
	}

	// Sort ProcessTypes alphabetically by process name
	sort.Sort(pts)

	return pts
}
