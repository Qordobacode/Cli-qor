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
	"github.com/spf13/cobra"
	"strings"
	"sync"
)

const (
	original = "original"
)

// downloadCmd represents the download command
var (
	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Downloads selected files",
		Long:  "Default file download command will give you two things  A)only the completed files B) will give you all the files (all locals and audiences without source file)",
		Run:   pullCommand,
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

func pullCommand(cmd *cobra.Command, args []string) {
	config, err := general.LoadConfig()
	if err != nil {
		return
	}
	workspace, err := general.GetWorkspace(config)
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	audiences := config.GetAudiences()
	if downloadAudience != "" {
		audienceList := strings.Split(downloadAudience, ",")
		audiences = make(map[string]bool)
		for _, lang := range audienceList {
			audiences[lang] = true
		}
	}
	for _, persona := range workspace.TargetPersonas {
		if _, ok := audiences[persona.Code]; len(audiences) > 0 && !ok {
			continue
		}
		files, err := general.GetFilesForTargetPerson(config, persona.ID)
		if err != nil {
			continue
		}
		wg.Add(len(files))
		for i := range files {
			go handleFile(config, persona.ID, &files[i], &wg)
		}
	}
	wg.Wait()
}

func handleFile(config *general.Config, personaID int, file *general.File, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	if !file.Completed && !isDownloadCurrent {
		// isDownloadCurrent - skip files with version
		return
	}
	if isDownloadSource || isDownloadOriginal {
		if isDownloadSource {
			fileName := general.BuildFileName(file, "")
			general.DownloadSourceFile(config, fileName, file, true)
		}
		if isDownloadOriginal {
			suffix := ""
			if isDownloadSource {
				// note if the customer using -s and -o in the same command rename the file original to filename-original.xxx
				suffix = original
			}
			fileName := general.BuildFileName(file, suffix)
			general.DownloadSourceFile(config, fileName, file, false)
		}
	} else {
		fileName := general.BuildFileName(file, "")
		if !isPullSkip || !general.FileExists(fileName) {
			general.DownloadFile(config, personaID, fileName, file)
		}
	}
}
