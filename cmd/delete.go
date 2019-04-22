package cmd

import (
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"
	"github.com/spf13/cobra"
)

var (
	deleteFileVersion string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete files from workspace",
	Run:   deleteFile,
}

func deleteFile(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		qordobaConfig, err := general.LoadConfig()
		if err != nil {
			return
		}
		general.FindFileAndDelete(qordobaConfig, args[0], deleteFileVersion)
	} else {
		log.Infof("No files to delete were specified")
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteFileVersion, "version", "", "version of file to delete")
}
