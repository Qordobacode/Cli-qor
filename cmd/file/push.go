package file

import (
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/spf13/cobra"
	"path/filepath"
)

var (
	pushVersion string
	files       string
	folderPath  string
)

// NewPushCmd creates `push` command
func NewPushCmd() *cobra.Command {
	// pushCmd represents the push command
	pushCmd := &cobra.Command{
		Annotations: map[string]string{"group": "file"},
		Use:         "push",
		Short:       "Push files or folders",
		Example:     `qor push --files testing.json --version 1.1 --verbose`,
		PreRun:      startLocalServices,
		Run:         pushCommand,
	}
	pushCmd.Flags().StringVarP(&pushVersion, "version", "v", "", "Set version to pushed file")
	pushCmd.Flags().StringVarP(&files, "files", "f", "", "Lists the file paths to upload")
	pushCmd.Flags().StringVarP(&folderPath, "file-path", "p", ".", "Push entire (relative) file paths")
	return pushCmd
}

func pushCommand(cmd *cobra.Command, args []string) {
	if appConfig == nil {
		log.Errorf("error occurred on configuration load")
		return
	}
	if folderPath == "" && files == "" && len(args) == 0 {
		pushSources := appConfig.Push.Sources

		log.Infof("no '--files' or '--file-path' params in command. 'source' param from config is used\n  File: %v\n  Folders: %v", pushSources.Files, pushSources.Folders)
		fileService.PushFiles(pushSources.Files, pushVersion)
		for _, folder := range pushSources.Folders {
			fileService.PushFolder(folder, pushVersion)
		}
		return
	}
	if files != "" || len(args) != 0 {
		fileList := filepath.SplitList(files)
		for _, arg := range args {
			argFiles := filepath.SplitList(arg)
			fileList = append(fileList, argFiles...)
		}
		for _, file := range fileList {
			fileService.PushFolder(file, pushVersion)
		}
	} else if folderPath != "" {
		log.Infof("Push folder %v", folderPath)
		fileService.PushFolder(folderPath, pushVersion)
	}
}
