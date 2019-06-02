package cmd

import (
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/spf13/cobra"
)

var (
	deleteFileVersion string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete files from workspace",
	PreRun: startLocalServices,
	Run:   deleteFile,
}

func deleteFile(cmd *cobra.Command, args []string) {
	if Config == nil {
		log.Errorf("error occurred on configuration load")
		return
	}
	if len(args) > 0 {
		FileService.DeleteFile(args[0], deleteFileVersion)
	} else {
		log.Infof("No files to delete were specified")
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteFileVersion, "version", "", "version of file to delete")
}
