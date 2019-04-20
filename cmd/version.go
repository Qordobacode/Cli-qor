package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	ApplicationName = "Qordoba CLI"
	ApiVersion      = "Qordoba Cli v4.0"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the Qordoba CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	fmt.Println(ApiVersion)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
