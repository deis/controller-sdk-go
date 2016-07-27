package api

// App is the definition of the app object.
type App struct {
	Created string `json:"created"`
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Updated string `json:"updated"`
	URL     string `json:"-"`
	UUID    string `json:"uuid"`
}

type Apps []App

func (a Apps) Len() int           { return len(a) }
func (a Apps) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Apps) Less(i, j int) bool { return a[i].ID < a[j].ID }

// AppCreateRequest is the definition of POST /v2/apps/.
type AppCreateRequest struct {
	ID string `json:"id,omitempty"`
}

// AppUpdateRequest is the definition of POST /v2/apps/<app id>/.
type AppUpdateRequest struct {
	Owner string `json:"owner,omitempty"`
}

// AppRunRequest is the definition of POST /v2/apps/<app id>/run.
type AppRunRequest struct {
	Command string `json:"command"`
}

// AppRunResponse is the definition of /v2/apps/<app id>/run.
type AppRunResponse struct {
	Output     string `json:"output"`
	ReturnCode int    `json:"exit_code"`
}
