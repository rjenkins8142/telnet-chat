package version

import "fmt"

var (
	// BuildDate is set automatically in the Makefile and contains the date the
	// build started.
	BuildDate string

	// GitBuild is set automatically in the Makefile and contains the last git
	// commit SHA.
	GitBuild string
)

// Info returns information about the current version/build.
func Info() string {
	if BuildDate != "" {
		return fmt.Sprintf("built %s [%s]", BuildDate, GitBuild)
	}
	return "v.unoffical"
}
