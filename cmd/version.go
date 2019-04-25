package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version-related constants
const (
	APIVersion = "Qordoba Cli v4.0"
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
	fmt.Println(APIVersion)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
