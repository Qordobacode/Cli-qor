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
	"sync"
)

// downloadCmd represents the download command
var (
	downloadCmd = &cobra.Command{
		Use:   "pull",
		Short: "Default file download command will give you two things  A)only the completed files B) will give you all the files (all locals and audiences without source file)",
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
	qordobaConfig, err := general.LoadConfig()
	if err != nil {
		return
	}
	workspace, err := general.GetWorkspace(qordobaConfig)
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	for _, persona := range workspace.TargetPersonas {
		files, err := general.GetFilesInWorkspace(qordobaConfig, persona.ID)
		if err != nil {
			continue
		}
		wg.Add(len(files))
		for _, file := range files {
			go general.DownloadFile(qordobaConfig, persona.ID, &file, &wg)
		}
	}
	wg.Wait()
}
