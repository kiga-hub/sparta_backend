package main

import (
	"fmt"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/kiga-hub/sparta_backend/cmd"
	_ "go.uber.org/automaxprocs"
)

var (
	// AppName - app name
	AppName string
	// AppVersion - app version
	AppVersion string
	// BuildVersion - build version
	BuildVersion string
	// BuildTime - build time
	BuildTime string
	// GitRevision - Git version
	GitRevision string
	// GitBranch - Git branch
	GitBranch string
	// GoVersion - Golang information
	GoVersion string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 2)
	spew.Config = *spew.NewDefaultConfig()
	spew.Config.ContinueOnMethod = true
	cmd.AppName = AppName
	cmd.AppVersion = AppVersion
	Version()

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

// Version prints version info of the program
func Version() {
	fmt.Printf(
		"App Name:\t%s\nApp Version:\t%s\nBuild version:\t%s\nBuild time:\t%s\nGit revision:\t%s\nGit branch:\t%s\nGolang Version: %s\n",
		AppName,
		AppVersion,
		BuildVersion,
		BuildTime,
		GitRevision,
		GitBranch,
		GoVersion,
	)
}
