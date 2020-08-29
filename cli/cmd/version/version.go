package version

import "fmt"

// populated by goreleaser
var (
	semver string
	commit string
	date   string
)

func String() string {
	return fmt.Sprintf("Version: %s, Commit: %s, Build Date: %s", semver, commit, date)
}

func GetCurrent() string {
	return semver
}
