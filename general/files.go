package general

import (
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/log"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	fileListURLTemplate                      = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files?withProgressStatus=%v"
	fileListURLTemplateWithLimit             = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files?withProgressStatus=%v&limit=%v"
	fileSearchURLTemplate                    = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files?withProgressStatus=%v&filename=%v&version=%v"
	fileDownloadTemplate                     = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files/%d/download"
	sourceFileDownloadTemplate               = "%s/v3/organizations/%d/workspaces/%d/files/%d/download/source?withUpdates=%v"
	fileDeleteTemplate                       = "%s/v3/organizations/%d/workspaces/%d/files/%d"
	defaultFilePerm              os.FileMode = 0666
)

var (
	forbiddenInFileNameSymbols, _ = regexp.Compile(`[:?!\\*/|<>]`)
)

// SearchForFiles function retrieves all files in workspace
func SearchForFiles(config *Config, personaID int, withProgressStatus bool) (*FileSearchResponse, error) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "SearchForFiles "+strconv.Itoa(personaID))
	}()
	base := config.GetAPIBase()
	fileListURL := fmt.Sprintf(fileListURLTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, personaID, withProgressStatus)
	return callFileRequestAndHandle(config, fileListURL)
}

func SearchForFilesWithLimit(config *Config, personaID int, withProgressStatus bool, limit int) (*FileSearchResponse, error) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "SearchForFilesWithLimit "+strconv.Itoa(personaID))
	}()
	base := config.GetAPIBase()
	fileListURL := fmt.Sprintf(fileListURLTemplateWithLimit, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, personaID, withProgressStatus, limit)
	return callFileRequestAndHandle(config, fileListURL)
}

func callFileRequestAndHandle(config *Config, getUserFiles string) (*FileSearchResponse, error) {
	fileBytesResponse, err := GetFromServer(config, getUserFiles)
	if err != nil {
		return nil, err
	}
	var response FileSearchResponse
	err = json.Unmarshal(fileBytesResponse, &response)
	if err != nil {
		log.Errorf("error occurred on server response unmarshalling: %v", err)
		return nil, err
	}
	return &response, nil
}

// DownloadFile function retrieves file in workspace
func DownloadFile(config *Config, personaID int, fileName string, file *File) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "DownloadFile")
	}()
	base := config.GetAPIBase()
	getFileContentURL := fmt.Sprintf(fileDownloadTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, personaID, file.FileID)
	fileBytesResponse, err := GetFromServer(config, getFileContentURL)
	if err != nil {
		log.Errorf("error occurred on file '%s' download: %v (url = %v)", fileName, err, getFileContentURL)
		return
	}
	log.Infof("file '%v' was downloaded", fileName)
	err = ioutil.WriteFile(fileName, fileBytesResponse, defaultFilePerm)
	if err != nil {
		log.Errorf("error occurred on writing file: %v", err)
	}
}

// DownloadSourceFile function retrieves all source files in workspace
func DownloadSourceFile(config *Config, fileName string, file *File, withUpdates bool) {
	base := config.GetAPIBase()
	getFileContentURL := fmt.Sprintf(sourceFileDownloadTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, file.FileID, withUpdates)
	fileBytesResponse, err := GetFromServer(config, getFileContentURL)
	if err != nil {
		log.Errorf("error occurred on file '%s' download: %v (url = %v)", fileName, err, getFileContentURL)
		return
	}
	log.Infof("source file '%v' was downloaded", fileName)
	err = ioutil.WriteFile(fileName, fileBytesResponse, defaultFilePerm)
	if err != nil {
		log.Errorf("error occurred on writing file: %v", err)
	}
}

// BuildFileName according to stored file name and version
func BuildFileName(file *File, suffix string) string {
	fileNames := strings.SplitN(file.Filename, ".", 2)
	if file.Version != "" {
		if suffix != "" {
			suffix = suffix + "_" + file.Version
		} else {
			suffix = file.Version
		}
	}
	resultName := file.Filename
	if suffix != "" {
		if len(fileNames) > 1 {
			resultName = fileNames[0] + "_" + suffix + "." + fileNames[1]
		}
		resultName = file.Filename + "_" + suffix
	}
	resultName = forbiddenInFileNameSymbols.ReplaceAllString(resultName, "")
	return resultName
}

// FindFile function
func FindFile(config *Config, fileName, version string) (*File, int) {
	log.Debugf("FindFile was called for file '%v'('%v')", fileName, version)
	workspace, err := GetWorkspace(config)
	if err != nil {
		return nil, 0
	}
	base := config.GetAPIBase()
	for _, persona := range workspace.TargetPersonas {
		fileListURL := fmt.Sprintf(fileSearchURLTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, persona.ID, false, fileName, version)
		fileSearchResponse, err := callFileRequestAndHandle(config, fileListURL)
		if err != nil {
			continue
		}
		for _, file := range fileSearchResponse.Files {
			if file.Filename == fileName {
				if file.Version == version {
					return &file, persona.ID
				}
			}
		}
	}
	if version == "" {
		log.Errorf("File '%s' WAS NOT FOUND", fileName)
	} else {
		log.Errorf("File '%s' with version '%s' WAS NOT FOUND", fileName, version)
	}
	return nil, 0
}

// FindFileAndDelete function retrieve file and delete it remotedly
func FindFileAndDelete(config *Config, fileName, version string) {
	log.Debugf("FindFileAndDelete was called for file '%v'('%v')", fileName, version)
	file, _ := FindFile(config, fileName, version)
	if file != nil {
		DeleteFile(config, file)
	}
}

// DeleteFile func delete file from parameters
func DeleteFile(config *Config, file *File) {
	base := config.GetAPIBase()
	deleteFileURL := fmt.Sprintf(fileDeleteTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, file.FileID)
	bytes, err := DeleteFromServer(config, deleteFileURL)
	if err != nil {
		return
	}
	var deleteResponse FileDeleteResponse
	err = json.Unmarshal(bytes, &deleteResponse)
	if err != nil {
		log.Errorf("error occurred on delete response unmarshalling: %v", err)
		return
	}
	if deleteResponse.Success {
		log.Infof("File '%s' with version '%s' was removed", file.Filename, file.Version)
	} else {
		log.Errorf("File '%s' with version '%s' WAS NOT REMOVED", file.Filename, file.Version)
	}
}
