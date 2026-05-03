package cmd

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/spf13/cobra"
)

var (
	Version  string
	Revision string
)

var rootCommand = &cobra.Command{
	Use: "grubzo",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		if c.Pprof {
			go func() { _ = http.ListenAndServe("0.0.0.0:6060", nil) }()
		}
	},
}

func init() {
	cobra.OnInitialize(func() {
		if err := loadConfig(); err != nil {
			log.Fatalf("failed to read config file: %v", err)
		}
	})

	rootCommand.AddCommand(
		serveCommand(),
		migrateCommand(),
	)
}

func Execute() error {
	return rootCommand.Execute()
}
