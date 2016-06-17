package deis

func checkAPICompatibility(serverAPIVersion string) error {
	if serverAPIVersion != APIVersion {
		return ErrAPIMismatch
	}

	return nil
}
