package general

import (
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/log"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	fileListURLTemplate                    = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files?withProgressStatus=%v"
	fileDownloadTemplate                   = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files/%d/download"
	sourceFileDownloadTemplate             = "%s/v3/organizations/%d/workspaces/%d/files/%d/download/source?withUpdates=%v"
	fileDeleteTemplate                     = "%s/v3/organizations/%d/workspaces/%d/files/%d"
	defaultFilePerm            os.FileMode = 0666
)

// SearchForFiles function retrieves all files in workspace
func SearchForFiles(config *Config, personaID int, withProgressStatus bool) (*FileSearchResponse, error) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "SearchForFiles "+strconv.Itoa(personaID))
	}()
	base := config.GetAPIBase()
	getUserFiles := fmt.Sprintf(fileListURLTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, personaID, withProgressStatus)
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
func DownloadFile(config *Config, personaID int, fileName string, file *Files) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "DownloadFile")
	}()
	base := config.GetAPIBase()
	getFileContent := fmt.Sprintf(fileDownloadTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, personaID, file.FileID)
	fileBytesResponse, err := GetFromServer(config, getFileContent)
	if err != nil {
		log.Errorf("error occurred on file %s download", fileName)
		return
	}
	log.Infof("file '%v' was downloaded", fileName)
	err = ioutil.WriteFile(fileName, fileBytesResponse, defaultFilePerm)
}

// DownloadSourceFile function retrieves all source files in workspace
func DownloadSourceFile(config *Config, fileName string, file *Files, withUpdates bool) {
	base := config.GetAPIBase()
	getFileContent := fmt.Sprintf(sourceFileDownloadTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, file.FileID, withUpdates)
	fileBytesResponse, err := GetFromServer(config, getFileContent)
	if err != nil {
		log.Errorf("error occurred on file %s download", fileName)
		return
	}
	log.Infof("source file '%v' was downloaded", fileName)
	err = ioutil.WriteFile(fileName, fileBytesResponse, defaultFilePerm)
}

// BuildFileName according to stored file name and version
func BuildFileName(file *Files, suffix string) string {
	fileNames := strings.SplitN(file.Filename, ".", 2)
	if file.Version != "" {
		if suffix != "" {
			suffix = suffix + "_" + file.Version
		} else {
			suffix = file.Version
		}
	}
	if suffix != "" {
		if len(fileNames) > 1 {
			return fileNames[0] + "_" + suffix + "." + fileNames[1]
		}
		return file.Filename + "_" + suffix
	}
	return file.Filename
}

// FindFile function
func FindFile(config *Config, fileName, version string) *Files {
	log.Debugf("FindFile was called for file '%v'('%v')", fileName, version)
	workspace, err := GetWorkspace(config)
	if err != nil {
		return nil
	}
	for _, persona := range workspace.TargetPersonas {
		fileSearchResponse, err := SearchForFiles(config, persona.ID, false)
		if err != nil {
			continue
		}
		for _, file := range fileSearchResponse.Files {
			if file.Filename == fileName {
				if file.Version == version {
					return &file
				}
			}
		}
	}
	log.Errorf("File '%s' with version '%s' WAS NOT FOUND", fileName, version)
	return nil
}

// FindFileAndDelete function retrieve file and delete it remotedly
func FindFileAndDelete(config *Config, fileName, version string) {
	log.Debugf("FindFileAndDelete was called for file '%v'('%v')", fileName, version)
	file := FindFile(config, fileName, version)
	if file != nil {
		DeleteFile(config, file)
	}
}

// DeleteFile func delete file from parameters
func DeleteFile(config *Config, file *Files) {
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
