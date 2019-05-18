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

package cmd

import (
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"
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
	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Downloads selected files",
		Long:  "Default file download command will give you two things  A)only the completed files B) will give you all the files (all locals and audiences without source file)",
		Run:   downloadCommand,
	}
	isDownloadCurrent  = false
	downloadAudience   = ""
	isDownloadSource   = false
	isDownloadOriginal = false
	isPullSkip         = false
)

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().BoolVarP(&isDownloadCurrent, "current", "c", false, "Pull the current state of the files")
	downloadCmd.Flags().StringVarP(&downloadAudience, "audience", "a", "", "Option to work only on specific (comma-separated) languages. example: `qor pull -a en-us,de-de`")
	downloadCmd.Flags().BoolVarP(&isDownloadSource, "source", "s", false, "File option to download the update source file")
	downloadCmd.Flags().BoolVarP(&isDownloadOriginal, "original", "o", false, " option to download the original file (note if the customer using -s and -o in the same command rename the file original to; filename-original.xxx) ")
	downloadCmd.Flags().BoolVar(&isPullSkip, "skip", false, "File option to download the update source file")
}

func downloadCommand(cmd *cobra.Command, args []string) {
	config, err := general.LoadConfig()
	if err != nil {
		return
	}
	workspace, err := general.GetWorkspace(config)
	if err != nil {
		return
	}
	files2Download := getFiles2Download(config, workspace)
	jobs := make(chan *File2Download, 1000)
	results := make(chan struct{}, 1000)

	for i := 0; i < 3; i++ {
		go worker(config, jobs, results)
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

func worker(config *general.Config, jobs chan *File2Download, results chan struct{}) {
	for j := range jobs {
		handleFile(config, j.PersonaID, j.File)
		results <- struct{}{}
	}
}

type File2Download struct {
	File      *general.File
	PersonaID int
}

func getFiles2Download(config *general.Config, workspace *general.Workspace) []*File2Download {
	audiences := config.GetAudiences()
	if downloadAudience != "" {
		audienceList := strings.Split(downloadAudience, ",")
		audiences = make(map[string]bool)
		for _, lang := range audienceList {
			audiences[lang] = true
		}
	}
	files2Download := make([]*File2Download, 0)
	for _, persona := range workspace.TargetPersonas {
		if _, ok := audiences[persona.Code]; len(audiences) > 0 && !ok {
			continue
		}
		response, err := general.SearchForFiles(config, persona.ID, false)
		if err != nil {
			continue
		}
		files := response.Files
		for i := range files {
			files2Download = append(files2Download, &File2Download{
				File:      &files[i],
				PersonaID: persona.ID,
			})
		}
	}
	return files2Download
}

func handleFile(config *general.Config, personaID int, file *general.File) {
	if !file.Completed && !isDownloadCurrent {
		// isDownloadCurrent - skip files with version
		return
	}
	if file.ErrorID != 0 || !file.Enabled {
		handleInvalidFile(file)
		return
	}
	if isDownloadSource || isDownloadOriginal {
		if isDownloadSource {
			downloadSourceFile(file, config)
		}
		if isDownloadOriginal {
			downloadOriginalFile(file, config)
		}
	} else {
		downloadFile(file, config, personaID)
	}
}

func handleInvalidFile(file *general.File) {
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

func downloadFile(file *general.File, config *general.Config, personaID int) {
	fileName := general.BuildFileName(file, "")
	if !isPullSkip || !general.FileExists(fileName) {
		general.DownloadFile(config, personaID, fileName, file)
		atomic.AddUint64(&ops, 1)
	}
}

func downloadSourceFile(file *general.File, config *general.Config) {
	fileName := general.BuildFileName(file, "")
	general.DownloadSourceFile(config, fileName, file, true)
	atomic.AddUint64(&ops, 1)
}

func downloadOriginalFile(file *general.File, config *general.Config) {
	suffix := ""
	if isDownloadSource {
		// note if the customer using -s and -o in the same command rename the file original to filename-original.xxx
		suffix = original
	}
	fileName := general.BuildFileName(file, suffix)
	general.DownloadSourceFile(config, fileName, file, false)
	atomic.AddUint64(&ops, 1)
}
