package api

import (
	"strings"
	"testing"
)

func TestAutoscaleString(t *testing.T) {
	a := Autoscale{}

	expected := strings.TrimSpace(`Min Replicas: 0
Max Replicas: 0
CPU: 0%`)

	if strings.TrimSpace(a.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, a.String())
	}

	a2 := Autoscale{
		Min:        3,
		Max:        8,
		CPUPercent: 40,
	}

	expected2 := strings.TrimSpace(`Min Replicas: 3
Max Replicas: 8
CPU: 40%`)

	if strings.TrimSpace(a2.String()) != expected2 {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected2, a2.String())
	}
}
