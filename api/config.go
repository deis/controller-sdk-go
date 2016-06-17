package api

import (
	"bytes"
	"fmt"
	"text/template"
)

// ConfigSet is the definition of POST /v2/apps/<app id>/config/.
type ConfigSet struct {
	Values map[string]string `json:"values"`
}

// ConfigUnset is the definition of POST /v2/apps/<app id>/config/.
type ConfigUnset struct {
	Values map[string]interface{} `json:"values"`
}

// Config is the structure of an app's config.
type Config struct {
	Owner       string                  `json:"owner,omitempty"`
	App         string                  `json:"app,omitempty"`
	Values      map[string]interface{}  `json:"values,omitempty"`
	Memory      map[string]interface{}  `json:"memory,omitempty"`
	CPU         map[string]interface{}  `json:"cpu,omitempty"`
	Healthcheck map[string]*Healthcheck `json:"healthcheck,omitempty"`
	Tags        map[string]interface{}  `json:"tags,omitempty"`
	Registry    map[string]interface{}  `json:"registry,omitempty"`
	Created     string                  `json:"created,omitempty"`
	Updated     string                  `json:"updated,omitempty"`
	UUID        string                  `json:"uuid,omitempty"`
}

// Healthcheck is the structure for an application healthcheck.
// Healthchecks only need to provide information about themselves.
// All the information is pushed to the server and handled by kubernetes.
type Healthcheck struct {
	InitialDelaySeconds int             `json:"initialDelaySeconds"`
	TimeoutSeconds      int             `json:"timeoutSeconds"`
	PeriodSeconds       int             `json:"periodSeconds"`
	SuccessThreshold    int             `json:"successThreshold"`
	FailureThreshold    int             `json:"failureThreshold"`
	Exec                *ExecProbe      `json:"exec,omitempty"`
	HTTPGet             *HTTPGetProbe   `json:"httpGet,omitempty"`
	TCPSocket           *TCPSocketProbe `json:"tcpSocket,omitempty"`
}

// String displays the HealthcheckHTTPGetProbe in a readable format.
func (h Healthcheck) String() string {
	var doc bytes.Buffer
	tmpl, err := template.New("healthcheck").Parse(`Initial Delay (seconds): {{.InitialDelaySeconds}}
Timeout (seconds): {{.TimeoutSeconds}}
Period (seconds): {{.PeriodSeconds}}
Success Threshold: {{.SuccessThreshold}}
Failure Threshold: {{.FailureThreshold}}
Exec Probe: {{or .Exec "N/A"}}
HTTP GET Probe: {{or .HTTPGet "N/A"}}
TCP Socket Probe: {{or .TCPSocket "N/A"}}`)
	if err != nil { panic(err) }
	err = tmpl.Execute(&doc, h)
	if err != nil { panic(err) }
	return doc.String()
}

// KVPair is a key/value pair used to parse values from
// strings into a formal structure.
type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (k KVPair) String() string {
	return k.Key+"="+k.Value
}

// ExecProbe executes a command within a Pod.
type ExecProbe struct {
	Command []string `json:"command"`
}

// String displays the ExecProbe in a readable format.
func (e ExecProbe) String() string {
	return fmt.Sprintf(`Command=%s`, e.Command)
}

// HTTPGetProbe performs an HTTP GET request to the Pod
// with the given path, port and headers.
type HTTPGetProbe struct {
	Path        string    `json:"path,omitempty"`
	Port        int       `json:"port"`
	HTTPHeaders []*KVPair `json:"httpHeaders,omitempty"`
}

// String displays the HTTPGetProbe in a readable format.
func (h HTTPGetProbe) String() string {
	return fmt.Sprintf(`Path="%s" Port=%d HTTPHeaders=%s`,
		h.Path,
		h.Port,
		h.HTTPHeaders)
}

// TCPSocketProbe attempts to open a socket connection to the
// Pod on the given port.
type TCPSocketProbe struct {
	Port int `json:"port"`
}

// String displays the TCPSocketProbe in a readable format.
func (t TCPSocketProbe) String() string {
	return fmt.Sprintf("Port=%d", t.Port)
}
