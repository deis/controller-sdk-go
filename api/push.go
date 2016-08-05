package api

// PushRequest is used to define the structure of the push hook
type PushRequest struct {
	Sha         string `json:"sha"`
	User        string `json:"receive_user"`
	App         string `json:"receive_repo"`
	Fingerprint string `json:"fingerprint"`
	Connection  string `json:"ssh_connection"`
	Command     string `json:"ssh_original_command"`
}
