package deis

import "strings"

func checkAPICompatibility(serverAPIVersion, clientAPIVersion string) error {
	sVersion := strings.Split(serverAPIVersion, ".")
	aVersion := strings.Split(clientAPIVersion, ".")

	// If API Versions are invalid, return a mismatch error.
	if len(sVersion) < 2 || len(aVersion) < 2 {
		return ErrAPIMismatch
	}

	// If major versions are different, return a mismatch error.
	if sVersion[0] != aVersion[0] {
		return ErrAPIMismatch
	}

	// If server is older than client, return mismatch error.
	if sVersion[1] < aVersion[1] {
		return ErrAPIMismatch
	}

	return nil
}
