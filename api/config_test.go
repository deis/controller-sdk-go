package api

import (
	"strings"
	"testing"
)

func TestHealthcheckString(t *testing.T) {
	h := Healthcheck{}

	expected := strings.TrimSpace(`Initial Delay (seconds): 0
Timeout (seconds): 0
Period (seconds): 0
Success Threshold: 0
Failure Threshold: 0
Exec Probe: N/A
HTTP GET Probe: N/A
TCP Socket Probe: N/A`)

	if strings.TrimSpace(h.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, h.String())
	}

	h.HTTPGet = &HTTPGetProbe{
		Path:        "/",
		Port:        80,
		HTTPHeaders: []*KVPair{{Name: "X-DEIS-IS", Value: "AWESOME"}},
	}

	expected = strings.TrimSpace(`Initial Delay (seconds): 0
Timeout (seconds): 0
Period (seconds): 0
Success Threshold: 0
Failure Threshold: 0
Exec Probe: N/A
HTTP GET Probe: Path="/" Port=80 HTTPHeaders=[X-DEIS-IS=AWESOME]
TCP Socket Probe: N/A`)

	if strings.TrimSpace(h.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, h.String())
	}

	h.Exec = &ExecProbe{Command: []string{"echo", "hi"}}

	h.TCPSocket = &TCPSocketProbe{
		Port: 80,
	}

	expected = strings.TrimSpace(`Initial Delay (seconds): 0
Timeout (seconds): 0
Period (seconds): 0
Success Threshold: 0
Failure Threshold: 0
Exec Probe: Command=[echo hi]
HTTP GET Probe: Path="/" Port=80 HTTPHeaders=[X-DEIS-IS=AWESOME]
TCP Socket Probe: Port=80`)

	if strings.TrimSpace(h.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, h.String())
	}
}
