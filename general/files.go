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
func DownloadFile(config *Config, personaID int, file *File, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	base := config.GetAPIBase()
	getFileContent := fmt.Sprintf(fileDownloadTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, personaID, file.FileID)
	fileBytesResponse, err := GetFromServer(config, getFileContent)
	if err != nil {
		log.Errorf("error occurred on file download: %v", err)
		return
	}
	fileName := BuildFileName(file)
	err = ioutil.WriteFile(fileName, fileBytesResponse, defaultFilePerm)
}

func BuildFileName(file *File) string {
	if file.Filename == "" {
		return file.Version
	}
	fileNames := strings.SplitN(file.Filename, ".", 2)
	if len(fileNames) > 1 {
		return fileNames[0] + "-" + file.Version + "." + fileNames[1]
	}
	return file.Filename + file.Version
}
