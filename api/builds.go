package api

// Build is the structure of the build object.
type Build struct {
	App        string            `json:"app"`
	Created    string            `json:"created"`
	Dockerfile string            `json:"dockerfile,omitempty"`
	Image      string            `json:"image,omitempty"`
	Owner      string            `json:"owner"`
	Procfile   map[string]string `json:"procfile"`
	Sha        string            `json:"sha,omitempty"`
	Updated    string            `json:"updated"`
	UUID       string            `json:"uuid"`
}

// CreateBuildRequest is the structure of POST /v2/apps/<app id>/builds/.
type CreateBuildRequest struct {
	Image    string            `json:"image"`
	Procfile map[string]string `json:"procfile,omitempty"`
}

// BuildHookRequest is a hook request to create a new build.
type BuildHookRequest struct {
	Sha        string      `json:"sha"`
	User       string      `json:"receive_user"`
	App        string      `json:"receive_repo"`
	Image      string      `json:"image"`
	Procfile   ProcessType `json:"procfile"`
	Dockerfile string      `json:"dockerfile"`
}
