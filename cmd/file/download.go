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
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

const (
	original           = "original"
	languageCode       = "language_code"
	languageLangCode   = "language_lang_code"
	languageName       = "language_name"
	languageNameCap    = "language_name_cap"
	localCapitalized   = "local_capitalized"
	languageNameAllcap = "language_name_allcap"
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
	isDownloadSkip     = false
	filePathPattern    = ""
	isFilePathPattern  = false
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
	downloadCmd.Flags().BoolVar(&isDownloadSkip, "skip", false, "File option to download the update source file")
	downloadCmd.Flags().StringVar(&filePathPattern, "file-path-pattern", "",
		`Download all target languages, or use in combination with -a flag. Replaces language pattern in path using provided variant:
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
	if !isFilePathPatternValid() {
		log.Errorf(`Invalid file-path-pattern "%s"; please provide one of:
- language_code 
- language_lang_code
- language_name
- language_name_cap
- local_capitalized
- language_name_allcap`, filePathPattern)
		return
	}

	workspace, err := workspaceService.LoadWorkspace()
	if err != nil || workspace == nil {
		return
	}
	isSourceCode, wasFound := validateWorkspace(workspace)
	if isSourceCode {
		return
	}
	// there might be a chance that in cached workspace we missed new language. Reload workspace from server then.
	if !wasFound {
		workspace, err := workspaceService.WorkspaceFromServer()
		if err != nil || workspace == nil {
			return
		}
		_, wasFound := validateWorkspace(workspace)
		if !wasFound {
			return
		}
	}
	if filePathPattern == "" && appConfig.Download.Target == "" && !isDownloadOriginal {
		log.Infof("Please update configuration and set the `download.target` field. For example `<language_code>-<filename>.<extension>`")
		os.Exit(1)
	}
	if isDownloadCurrent && isDownloadOriginal {
		log.Infof("-c parameter has no effect when used with -o, proceeding with downloading the original version of file.")
	}
	isFilePathPattern = filePathPattern != ""
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

func validateWorkspace(workspace *types.WorkspaceData) (isSourceCode bool, wasFound bool) {
	if downloadAudience != "" {
		// allow audiences like `qor download -a ja-jp,ko-kr`
		audienceList := strings.Split(downloadAudience, ",")
	A:
		for _, audience := range audienceList {
			targetCodesSlice := make([]string, 0)
			for _, t := range workspace.Workspace.TargetPersonas {
				if t.Code == audience {
					continue A
				}
				targetCodesSlice = append(targetCodesSlice, "`"+t.Code+"`")
			}
			targetLanguages := strings.Join(targetCodesSlice, ", ")
			if workspace.Workspace.SourcePersona.Code == audience {
				log.Errorf("`%s` is a source language. Please use the -o and -s parameters to download source files.", audience)
				return true, false
			}

			log.Errorf("`%s` does not match one of available project target languages: %s", audience, targetLanguages)
			return false, false
		}
	}
	return false, true
}

func isFilePathPatternValid() bool {
	return filePathPattern == "" || filePathPattern == languageCode || filePathPattern == languageLangCode ||
		filePathPattern == languageName || filePathPattern == languageNameCap ||
		filePathPattern == languageNameAllcap || filePathPattern == localCapitalized
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

func buildReplaceInString(person types.Person, filePathPattern string) (string, map[string]string) {
	replacementMap := make(map[string]string)
	codes := strings.Split(person.Code, "-")
	names := strings.Split(person.Name, "-")
	if len(codes) < 2 || len(names) < 2 {
		return "", replacementMap
	}
	replacementMap["<language_code>"] = person.Code
	replacementMap["<language_lang_code>"] = strings.TrimSpace(codes[0])
	replacementMap["<language_name>"] = strings.ToLower(strings.TrimSpace(names[0]))
	replacementMap["<language_name_cap>"] = strings.Title(strings.TrimSpace(names[0]))
	replacementMap["<language_name_allcap>"] = strings.ToUpper(strings.TrimSpace(names[0]))
	replacementMap["<local_capitalized>"] = strings.ToUpper(strings.TrimSpace(codes[1]))
	if replacementMap[filePathPattern] == "" {
		return replacementMap["<language_lang_code>"], replacementMap
	}
	return replacementMap[filePathPattern], replacementMap
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
	for pi := range workspace.TargetPersonas {
		persona := workspace.TargetPersonas[pi]
		if _, ok := audiences[persona.Code]; len(audiences) > 0 && !ok {
			continue
		}
		response, err := fileService.WorkspaceFiles(persona.ID, false)
		if err != nil {
			continue
		}
		files := response.Files
		for i := range files {
			replaceIn, replaceMap := buildReplaceInString(persona, filePathTemplate)
			files2Download = append(files2Download, &types.File2Download{
				File:       &files[i],
				Person:     workspace.TargetPersonas[pi],
				ReplaceIn:  replaceIn,
				ReplaceMap: replaceMap,
			})
		}
		if isDownloadOriginal || isDownloadSource {
			break
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
	if isDownloadSource && !(filePathPattern == "" && appConfig.Download.Target == "") {
		downloadSourceFile(j)
	}
	if isDownloadOriginal {
		downloadOriginalFile(j, matchFilepathName)
	}
	if !isDownloadOriginal && !isDownloadSource {
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
	dir := filepath.Dir(j.File.Filepath)
	if appConfig.Download.Target != "" && dir != "" && dir != "." {
		log.Infof("[TARGET] file '%s' has file path. File path is not supported with config`download.target`. Skip.", j.File.Filepath)
		return
	}
	fileName := local.BuildDirectoryFilePath(j, matchFilepathName, "", isFilePathPattern)
	if !isDownloadSkip || !local.FileExists(fileName) {
		fileService.DownloadFile(j.Person, fileName, j.File)
		atomic.AddUint64(&ops, 1)
	}
}

func downloadSourceFile(j *types.File2Download) {
	dir := filepath.Dir(j.File.Filepath)
	if appConfig.Download.Target != "" && dir != "" && dir != "." {
		log.Infof("[SOURCE] file '%s' has file path. File path is not supported with config `download.target`. Skip.", j.File.Filepath)
		return
	}
	fileName := local.BuildDirectoryFilePath(j, []string{}, "", true)
	if !isDownloadSkip || !local.FileExists(fileName) {
		fileService.DownloadSourceFile(fileName, j.File, true)
		atomic.AddUint64(&ops, 1)
	}
}

func downloadOriginalFile(j *types.File2Download, matchFilepathName []string) {
	suffix := ""
	if isDownloadSource {
		// note if the customer using -s and -o in the same command rename the file original to filename-original.xxx
		suffix = original
	}
	fileName := local.BuildDirectoryFilePath(j, []string{}, suffix, true)
	if !isDownloadSkip || !local.FileExists(fileName) {
		fileService.DownloadSourceFile(fileName, j.File, false)
		atomic.AddUint64(&ops, 1)
	}
}
