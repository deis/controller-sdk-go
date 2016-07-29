package deis

import "testing"

type versionComparison struct {
	Client string
	Server string
	Error  error
}

func TestCheckAPIVersions(t *testing.T) {
	comparisons := []versionComparison{
		{"1.2", "2.1", ErrAPIMismatch},
		{"2.1", "1.2", ErrAPIMismatch},
		{"2.1", "2.2", ErrAPIMismatch},
		{"2.3", "2.0", nil},
	}

	for _, check := range comparisons {
		err := checkAPICompatibility(check.Client, check.Server)

		if err != check.Error {
			t.Errorf("%v: Expected %v, Got %v", check, check.Error, err)
		}
	}
}
