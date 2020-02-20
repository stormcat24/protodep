package helper

import (
	"os"
)

// IsAvailableSSH is Check whether this machine can use git protocol
func IsAvailableSSH(identifyPath string) (bool, error) {
	if _, err := os.Stat(identifyPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// TODO: validate ssh key
	return true, nil
}
