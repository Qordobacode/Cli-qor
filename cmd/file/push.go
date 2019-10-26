package file

import (
	"github.com/qordobacode/cli-v2/pkg/file"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	pushVersion string
	files       string
	isFilePath  bool
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
	pushCmd.Flags().BoolVarP(&isFilePath, "file-path", "p", false, "Reads push.sources.folders from config file and push its content to server")
	return pushCmd
}

func pushCommand(cmd *cobra.Command, args []string) {
	if appConfig == nil {
		log.Errorf("error occurred on configuration load")
		return
	}
	if isFilePath && appConfig.Download.Target != "" {
		log.Errorf("Please remove `download.target` from your configuration file; it is not supported with file paths.")
		return
	}

	if !isFilePath && files == "" && len(args) == 0 {
		pushSources := appConfig.Push.Sources

		log.Infof("no '--files' or '--file-path' params in command. 'push.source' param from config is used\n  File: %v\n  Folders: %v", pushSources.Files, pushSources.Folders)
		fileService.PushFiles(pushSources.Files, pushVersion, false)
		for _, folder := range pushSources.Folders {
			if !filepath.IsAbs(folder) {
				log.Errorf("Please provide an absolute path for config parameter `push.sources.folders`")
				os.Exit(1)
			}
			fileService.PushFolder(folder, pushVersion, isFilePath)
		}
		if file.TotalSkipped > 0 {
			log.Infof(`%v files were skipped as their extension did not match one of: %s`, file.TotalSkipped, file.MimeTypes)
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
			fileService.PushFolder(file, pushVersion, isFilePath)
		}
	} else if isFilePath {
		if len(appConfig.Push.Sources.Folders) > 0 {
			log.Infof("Files being recursively pushed from path provided in configuration at `push.sources.folders`")
			fileService.PushFolder(appConfig.Push.Sources.Folders[0], pushVersion, isFilePath)
		} else {
			log.Errorf("--file-path variants uses push.sources.folders from config and push it on server")
		}
	}
}
