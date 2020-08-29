package version

import "fmt"

// populated by goreleaser
var (
	semver string
	commit string
	date   string
)

func String() string {
	return fmt.Sprintf("kubectl flame version: %s\ncommit: %s\nbuild date: %s", semver, commit, date)
}

func GetCurrent() string {
	return semver
}
