package cmd

import (
	"fmt"
	"github.com/qordobacode/cli-v2/cmd/config"
	"github.com/qordobacode/cli-v2/cmd/file"
	"github.com/qordobacode/cli-v2/cmd/info"
	"github.com/qordobacode/cli-v2/cmd/segment"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:     "qor",
		Short:   "Qordoba CLI",
		Long:    `This CLI is used for simplified access to Qordoba API`,
		Version: info.APIVersion,
		Run: func(cmd *cobra.Command, args []string) {
			if Version {
				info.PrintVersion()
			} else {
				cmd.Help()
			}
		},
	}
	Version bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Version, "version", "v", false, "Get version of CLI")
	rootCmd.PersistentFlags().BoolVar(&log.IsVerbose, "verbose", false, "Print verbose output")

	rootCmd.AddCommand(
		config.NewInitCmd(),

		file.NewPushCmd(),
		file.NewDownloadCommand(),
		file.NewDeleteFileCmd(),

		segment.NewAddKeyCommand(),
		segment.NewUpdateSegmentCommand(),
		segment.NewDeleteSegmentCommand(),
		segment.NewValueKeyCommand(),

		info.NewCmdVersion(),
		info.NewLsCommand(),
		info.NewStatusCommand(),
		info.NewScoreCommand(),
	)
}
