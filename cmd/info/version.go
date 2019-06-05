package info

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version-related constants
const (
	APIVersion = "Qordoba Cli v4.0"
)

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

func PrintVersion() {
	fmt.Println(APIVersion)
}
