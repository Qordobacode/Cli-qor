package cmd

import (
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	pushVersion string
	files       string
	folderPath  string

	// pushCmd represents the push command
	pushCmd = &cobra.Command{
		Use:    "push",
		Short:  "Push files or folders",
		PreRun: startLocalServices,
		Run:    pushCommand,
	}
)

func init() {
	pushCmd.Flags().StringVarP(&pushVersion, "version", "v", "", "--version")
	pushCmd.Flags().StringVarP(&files, "files", "f", "", "--update")
	pushCmd.Flags().StringVarP(&folderPath, "file-path", "p", "", "--update")
	rootCmd.AddCommand(pushCmd)
}

func pushCommand(cmd *cobra.Command, args []string) {
	if Config == nil {
		log.Errorf("error occurred on configuration load")
		return
	}
	if folderPath == "" && files == "" && len(args) == 0 {
		pushSources := Config.Push.Sources

		log.Infof("no '--files' or '--file-path' params in command. 'source' param from config is used\n  File: %v\n  Folders: %v", pushSources.Files, pushSources.Folders)
		FileService.PushFiles(pushSources.Files, pushVersion)
		for _, folder := range pushSources.Folders {
			FileService.PushFolder(folder, pushVersion)
		}
		return
	}
	if files != "" || len(args) != 0 {
		fileList := filepath.SplitList(files)
		for _, arg := range args {
			argFiles := filepath.SplitList(arg)
			fileList = append(fileList, argFiles...)
		}
		log.Debugf("Result list of files from line is: %v", string(os.PathListSeparator), fileList)
		for _, file := range fileList {
			FileService.PushFolder(file, pushVersion)
		}
	} else if folderPath != "" {
		FileService.PushFolder(folderPath, pushVersion)
	}
}
