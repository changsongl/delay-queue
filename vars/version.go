package vars

import "fmt"

// build information
var (
	BuildProgram string
	BuildGitPath string
	BuildVersion string
	BuildTime    string
	GoVersion    string
)

func BuildInfo() string {
	return fmt.Sprintf(
		"[Build Info] \nProgram   : %s \nVersion   : %s \nGo Version: %s \nBuild Time: %s \nGithub    : %s\n",
		BuildProgram, BuildVersion, GoVersion, BuildTime, BuildGitPath)
}
