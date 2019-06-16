// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/spf13/cobra"
	"strings"
	"sync/atomic"
	"time"
)

const (
	original = "original"
)

var (
	ops uint64
)

// downloadCmd represents the download command
var (
	isDownloadCurrent  = false
	downloadAudience   = ""
	isDownloadSource   = false
	isDownloadOriginal = false
	isPullSkip         = false
)

// NewDownloadCommand command create `download` command
func NewDownloadCommand() *cobra.Command {
	downloadCmd := &cobra.Command{
		Annotations: map[string]string{"group": "file"},
		Use:         "download",
		Short:       "Downloads selected files",
		Long:        "Default file download command will give you two things  A)only the completed files B) will give you all the files (all locals and audiences without source file)",
		Example:     `qor download -a en-us,de-de`,
		PreRun:      startLocalServices,
		Run:         downloadCommand,
	}

	downloadCmd.Flags().BoolVarP(&isDownloadCurrent, "current", "c", false, "Pull the current state of the files")
	downloadCmd.Flags().StringVarP(&downloadAudience, "audience", "a", "", "Option to work only on specific (comma-separated) languages. example: `qor pull -a en-us,de-de`")
	downloadCmd.Flags().BoolVarP(&isDownloadSource, "source", "s", false, "File option to download the update source file")
	downloadCmd.Flags().BoolVarP(&isDownloadOriginal, "original", "o", false, "Option to download the original file (note if the customer using -s and -o in the same command rename the file original to; filename-original.xxx) ")
	downloadCmd.Flags().BoolVar(&isPullSkip, "skip", false, "File option to download the update source file")
	return downloadCmd
}

func downloadCommand(cmd *cobra.Command, args []string) {
	if appConfig == nil {
		log.Errorf("error occurred on configuration load")
		return
	}
	workspace, err := workspaceService.LoadWorkspace()
	if err != nil || workspace == nil {
		return
	}
	files2Download := files2Download(&workspace.Workspace)
	jobs := make(chan *file2Download, 1000)
	results := make(chan struct{}, 1000)

	for i := 0; i < 3; i++ {
		go worker(jobs, results)
	}
	for _, file2Download := range files2Download {
		jobs <- file2Download
	}
	close(jobs)
	for i := 0; i < len(files2Download); i++ {
		<-results
	}

	// let all error logs go before final messages
	time.Sleep(time.Second)
	if isDownloadCurrent {
		log.Infof("downloaded %v files", ops)
	} else {
		log.Infof("downloaded %v completed files", ops)
	}
}

func worker(jobs chan *file2Download, results chan struct{}) {
	for j := range jobs {
		handleFile(j.PersonaID, j.File)
		results <- struct{}{}
	}
}

// file2Download struct describe chunk of download work
type file2Download struct {
	File      *types.File
	PersonaID int
}

func files2Download(workspace *types.Workspace) []*file2Download {
	audiences := appConfig.Audiences()
	if downloadAudience != "" {
		audienceList := strings.Split(downloadAudience, ",")
		audiences = make(map[string]bool)
		for _, lang := range audienceList {
			audiences[lang] = true
		}
	}
	files2Download := make([]*file2Download, 0)
	for _, persona := range workspace.TargetPersonas {
		if _, ok := audiences[persona.Code]; len(audiences) > 0 && !ok {
			continue
		}
		response, err := fileService.WorkspaceFiles(persona.ID, false)
		if err != nil {
			continue
		}
		files := response.Files
		for i := range files {
			files2Download = append(files2Download, &file2Download{
				File:      &files[i],
				PersonaID: persona.ID,
			})
		}
	}
	return files2Download
}

func handleFile(personaID int, file *types.File) {
	if !file.Completed && !isDownloadCurrent {
		// isDownloadCurrent - skip files with version
		log.Infof("file %s is not completed. Use flag '-c' or '--current' to download even not completed files", file.Filename)
		return
	}
	if file.ErrorID != 0 || !file.Enabled {
		handleInvalidFile(file)
		return
	}
	if isDownloadSource || isDownloadOriginal {
		if isDownloadSource {
			downloadSourceFile(file)
		}
		if isDownloadOriginal {
			downloadOriginalFile(file)
		}
	} else {
		downloadFile(file, personaID)
	}
}

func handleInvalidFile(file *types.File) {
	if file.ErrorID != 0 {
		if file.Version != "" {
			log.Errorf("'%s'(version '%v') has error. Skip its download", file.Filename, file.Version)
		} else {
			log.Errorf("'%s' has error. Skip its download", file.Filename)
		}
		return
	}
	if !file.Enabled {
		if file.Version != "" {
			log.Errorf("File '%s' (version '%s') is disabled. Skip its download", file.Filename, file.Version)
		} else {
			log.Errorf("File '%s' is disabled. Skip its download", file.Filename)
		}
		return
	}
}

func downloadFile(file *types.File, personaID int) {
	fileName := local.BuildFileName(file, "")
	if !isPullSkip || !local.FileExists(fileName) {
		fileService.DownloadFile(personaID, fileName, file)
		atomic.AddUint64(&ops, 1)
	}
}

func downloadSourceFile(file *types.File) {
	fileName := local.BuildFileName(file, "")
	fileService.DownloadSourceFile(fileName, file, true)
	atomic.AddUint64(&ops, 1)
}

func downloadOriginalFile(file *types.File) {
	suffix := ""
	if isDownloadSource {
		// note if the customer using -s and -o in the same command rename the file original to filename-original.xxx
		suffix = original
	}
	fileName := local.BuildFileName(file, suffix)
	fileService.DownloadSourceFile(fileName, file, false)
	atomic.AddUint64(&ops, 1)
}
