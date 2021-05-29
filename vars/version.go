package vars

import "fmt"

// build information
var (
	BuildProgram string
	BuildGitPath string
	BuildVersion string
	BuildTime    string
)

func BuildInfo() string {
	return fmt.Sprintf(
		"[Build Info] \nProgram: %s \nVersion: %s \nBuild Time: %s \nGithub: %s\n",
		BuildProgram, BuildVersion, BuildTime, BuildGitPath)
}
