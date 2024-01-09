package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/davecgh/go-spew/spew"
	"github.com/kiga-hub/arc/micro"
	basicComponent "github.com/kiga-hub/arc/micro/component"
	"github.com/kiga-hub/arc/tracing"
	"github.com/spf13/cobra"

	"github.com/kiga-hub/websocket/pkg/component"
)

func init() {
	spew.Config = *spew.NewDefaultConfig()
	spew.Config.ContinueOnMethod = true
}

// serverCmd .
var serverCmd = &cobra.Command{
	Use:   "run",
	Short: "run websocket",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	// recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", "recover", r)
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	server, err := micro.NewServer(
		AppName,
		AppVersion,
		[]micro.IComponent{
			&basicComponent.LoggingComponent{},
			&tracing.Component{},
			&component.WebScoketComponent{},
		},
	)
	if err != nil {
		panic(err)
	}
	err = server.Init()
	if err != nil {
		panic(err)
	}
	err = server.Run()
	if err != nil {
		panic(err)
	}
}
