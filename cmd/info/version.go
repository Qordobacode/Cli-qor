package info

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version-related constants
const (
	APIVersion = "Qordoba Cli v4"
)

var (
	VersionFlag = "0.7.5"
)

// NewCmdVersion function create `version` command
func NewCmdVersion() *cobra.Command {
	return &cobra.Command{
		Annotations: map[string]string{"group": "info"},
		Use:         "version",
		Short:       "Print the Qordoba CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			PrintVersion()
		},
	}
}

// PrintVersion function print current version to stdout
func PrintVersion() {
	fmt.Printf("%s-%s\n", APIVersion, VersionFlag)
}
