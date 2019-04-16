package cmd

import (
	"fmt"
	"github.com/qordobacode/cli-v2/files"
	"github.com/qordobacode/cli-v2/models"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init configuration for Qordoba CLI from STDIN",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("init called")
		fileName := args[0]
		return RunInit(fileName)
	},
	Example:     "qor init",
	Annotations: map[string]string{"version": ApiVersion},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// RunConfigCreate creates a config with the given options.
func RunInit(fileName string) error {
	var config *models.QordobaConfig
	var err error
	if fileName != "" {
		config, err = files.ReadConfigInPath(fileName)
		if err != nil {
			return err
		}
	}

	files.PersistAppConfig(config)
	return nil
}
