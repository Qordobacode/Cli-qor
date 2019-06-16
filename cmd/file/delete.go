package file

import (
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/spf13/cobra"
)

var (
	deleteFileVersion string
)

// NewDeleteFileCmd function build `delete` command for cobra
func NewDeleteFileCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Annotations: map[string]string{"group": "file"},
		Use:         "delete",
		Short:       "Delete files from workspace",
		Example:     "qor delete file_name.doc --version 1",
		PreRun:      startLocalServices,
		Run:         deleteFile,
	}
	deleteCmd.Flags().StringVar(&deleteFileVersion, "version", "", "version of file to delete")
	return deleteCmd
}

func deleteFile(cmd *cobra.Command, args []string) {
	if appConfig == nil {
		log.Errorf("error occurred on configuration load")
		return
	}
	if len(args) > 0 {
		fileService.DeleteFile(args[0], deleteFileVersion)
	} else {
		log.Infof("No files to delete were specified")
	}
}
