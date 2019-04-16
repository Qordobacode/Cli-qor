package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	ApplicationName = "Qordoba CLI"
	ApiVersion      = "v1"
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
	fmt.Printf("%s version: %s\n", ApplicationName, ApiVersion)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
