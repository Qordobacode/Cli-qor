package general

import (
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/log"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

const (
	fileListURLTemplate              = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files"
	fileDownloadTemplate             = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files/%d/download"
	fileDeleteTemplate               = "%s/v3/organizations/%d/workspaces/%d/files/%d"
	defaultFilePerm      os.FileMode = 0666
)

// DownloadFile function retrieves all files in workspace
func GetFilesInWorkspace(config *Config, personaID int) ([]File, error) {
	base := config.GetAPIBase()
	getUserFiles := fmt.Sprintf(fileListURLTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, personaID)
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
	return response.Files, nil
}

// DownloadFile function retrieves all files in workspace
func DownloadFile(config *Config, personaID int, fileName string, file *File, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	base := config.GetAPIBase()
	getFileContent := fmt.Sprintf(fileDownloadTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, personaID, file.FileID)
	fileBytesResponse, err := GetFromServer(config, getFileContent)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(fileName, fileBytesResponse, defaultFilePerm)
}

// BuildFileName according to stored file name and version
func BuildFileName(file *File) string {
	if file.Filename == "" {
		return file.Version
	}
	fileNames := strings.SplitN(file.Filename, ".", 2)
	if file.Version != "" {
		if len(fileNames) > 1 {
			return fileNames[0] + "_" + file.Version + "." + fileNames[1]
		}
		return file.Filename + "_" + file.Version
	}
	return file.Filename
}

func FindFileAndDelete(config *Config, fileName, version string) {
	log.Debugf("FindFileAndDelete was called for file '%v'('%v')", fileName, version)
	workspace, err := GetWorkspace(config)
	if err != nil {
		return
	}
	for _, persona := range workspace.TargetPersonas {
		files, err := GetFilesInWorkspace(config, persona.ID)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.Filename == fileName {
				if file.Version == version {
					DeleteFile(config, &file)
					return
				}
			}
		}
		log.Errorf("File '%s' with version '%s' WAS NOT FOUND", fileName, version)
	}
}

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