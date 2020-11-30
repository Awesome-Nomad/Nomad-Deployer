package semver

import (
	msemver "github.com/Masterminds/semver"
)

// IsSemver check if input is a valid semantic version or not.
func IsSemver(input string) bool {
	_, err := msemver.NewVersion(input)
	return err != nil
}
