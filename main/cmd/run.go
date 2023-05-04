package cmd

import (
	"github.com/ichaly/yugong/core"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"path/filepath"
)

var configFile string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "version subcommand show git version info.",

	Run: func(cmd *cobra.Command, args []string) {
		if configFile == "" {
			configFile = filepath.Join("../conf", "dev.yml")
		}
		fx.New(
			core.Modules,
			fx.Provide(
				fx.Annotated{
					Name:   "configFile",
					Target: func() string { return configFile },
				},
			),
		).Run()
	},
}

func init() {
	runCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "start app with config file")
	rootCmd.AddCommand(runCmd)
}
