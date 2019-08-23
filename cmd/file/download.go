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
	filePathPattern    = ""
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
		Run:         downloadFiles,
	}

	downloadCmd.Flags().BoolVarP(&isDownloadCurrent, "current", "c", false, "Pull the current state of the files")
	downloadCmd.Flags().StringVarP(&downloadAudience, "audience", "a", "", "Option to work only on specific (comma-separated) languages. example: `qor pull -a en-us,de-de`")
	downloadCmd.Flags().BoolVarP(&isDownloadSource, "source", "s", false, "File option to download the update source file")
	downloadCmd.Flags().BoolVarP(&isDownloadOriginal, "original", "o", false, "Option to download the original file (note if the customer using -s and -o in the same command rename the file original to; filename-original.xxx) ")
	downloadCmd.Flags().BoolVar(&isPullSkip, "skip", false, "File option to download the update source file")
	downloadCmd.Flags().StringVar(&filePathPattern, "file-path-pattern", "language_lang_code", `Source code pattern. Possible variants:
- language_code 
- language_lang_code
- language_name
- language_name_cap
- language_name_allcap
- local_capitalized
`)
	return downloadCmd
}

func downloadFiles(cmd *cobra.Command, args []string) {
	if appConfig == nil {
		log.Errorf("error occurred on configuration load")
		return
	}
	workspace, err := workspaceService.LoadWorkspace()
	if err != nil || workspace == nil {
		return
	}
	log.Infof("File Path Pattern used is `%s`", filePathPattern)
	matchFilepathName := buildPatternName(workspace.Workspace.SourcePersona)
	files2Download := files2Download(&workspace.Workspace, filePathPattern)
	jobs := make(chan *types.File2Download, 1000)
	results := make(chan struct{}, 1000)

	for i := 0; i < 3; i++ {
		go worker(jobs, results, matchFilepathName)
	}
	for _, files2Download := range files2Download {
		jobs <- files2Download
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

// buildPatternName builds
func buildPatternName(person types.Person) []string {
	results := make([]string, 0, 0)
	results = updateByVariantSlice(person.Code, results)
	results = updateByVariantSlice(person.Name, results)
	return results
}

func updateByVariantSlice(variantString string, results []string) []string {
	variants := strings.Split(variantString, "-")
	for _, variant := range variants {
		trimmedVariant := strings.TrimSpace(variant)
		results = append(results, strings.ToLower(trimmedVariant))
		results = append(results, trimmedVariant)
		results = append(results, strings.ToUpper(trimmedVariant))
	}
	return results
}

func buildReplaceInString(person types.Person, filePathPattern string) string {
	codes := strings.Split(person.Code, "-")
	names := strings.Split(person.Name, "-")
	if len(codes) < 2 || len(names) < 2 {
		return ""
	}
	switch filePathPattern {
	case "language_code":
		return person.Code
	case "language_lang_code":
		return strings.TrimSpace(codes[0])
	case "language_name":
		return strings.ToLower(strings.TrimSpace(names[0]))
	case "language_name_cap":
		return strings.Title(strings.TrimSpace(names[0]))
	case "language_name_allcap":
		return strings.ToUpper(strings.TrimSpace(names[0]))
	case "local_capitalized":
		return strings.ToUpper(strings.TrimSpace(codes[1]))
	default:
		return strings.TrimSpace(codes[0])
	}
}

func worker(jobs chan *types.File2Download, results chan struct{}, matchFilepathName []string) {
	for j := range jobs {
		handleFile(j, matchFilepathName)
		results <- struct{}{}
	}
}

func files2Download(workspace *types.Workspace, filePathTemplate string) []*types.File2Download {
	audiences := appConfig.Audiences()
	if downloadAudience != "" {
		audienceList := strings.Split(downloadAudience, ",")
		audiences = make(map[string]bool)
		for _, lang := range audienceList {
			audiences[lang] = true
		}
	}
	files2Download := make([]*types.File2Download, 0)
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
			replaceIn := buildReplaceInString(persona, filePathTemplate)
			files2Download = append(files2Download, &types.File2Download{
				File:      &files[i],
				PersonaID: persona.ID,
				ReplaceIn: replaceIn,
			})
		}
	}
	return files2Download
}

func handleFile(j *types.File2Download, matchFilepathName []string) {
	if !j.File.Completed && !isDownloadCurrent && !isDownloadOriginal {
		// isDownloadCurrent - skip files with version
		log.Infof("file %s is not completed. Use flag '-c' or '--current' to download even not completed files", j.File.Filename)
		return
	}
	if j.File.ErrorID != 0 || !j.File.Enabled {
		handleInvalidFile(j.File)
		return
	}
	if isDownloadSource || isDownloadOriginal {
		if isDownloadSource {
			downloadSourceFile(j)
		}
		if isDownloadOriginal {
			downloadOriginalFile(j, matchFilepathName)
		}
	} else {
		downloadFile(j, matchFilepathName)
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

func downloadFile(j *types.File2Download, matchFilepathName []string) {
	fileName := local.BuildDirectoryFilePath(j, matchFilepathName, "")
	if !isPullSkip || !local.FileExists(fileName) {
		fileService.DownloadFile(j.PersonaID, fileName, j.File)
		atomic.AddUint64(&ops, 1)
	}
}

func downloadSourceFile(j *types.File2Download) {
	fileName := j.File.Filepath
	fileService.DownloadSourceFile(fileName, j.File, true)
	atomic.AddUint64(&ops, 1)
}

func downloadOriginalFile(j *types.File2Download, matchFilepathName []string) {
	suffix := ""
	if isDownloadSource {
		// note if the customer using -s and -o in the same command rename the file original to filename-original.xxx
		suffix = original
	}
	fileName := local.BuildDirectoryFilePath(j, matchFilepathName, suffix)
	fileService.DownloadSourceFile(fileName, j.File, false)
	atomic.AddUint64(&ops, 1)
}
