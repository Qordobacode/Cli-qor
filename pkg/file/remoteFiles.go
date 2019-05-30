package file

import (
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"strconv"
	"time"
)

const (
	fileListURLTemplate          = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files?withProgressStatus=%v"
	fileListURLTemplateWithLimit = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files?withProgressStatus=%v&limit=%v"
	fileSearchURLTemplate        = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files?withProgressStatus=%v&filename=%v&version=%v"
	fileDownloadTemplate         = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files/%d/download"
	sourceFileDownloadTemplate   = "%s/v3/organizations/%d/workspaces/%d/files/%d/download/source?withUpdates=%v"
	fileDeleteTemplate           = "%s/v3/organizations/%d/workspaces/%d/files/%d"
)

type FileService struct {
	Config           *types.Config
	QordobaClient    pkg.QordobaClient
	WorkspaceService pkg.WorkspaceService
	Local            pkg.Local
}

// WorkspaceFiles function retrieves all files in workspace
func (f *FileService) WorkspaceFiles(personaID int, withProgressStatus bool) (*types.FileSearchResponse, error) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "WorkspaceFiles "+strconv.Itoa(personaID))
	}()
	base := f.Config.GetAPIBase()
	fileListURL := fmt.Sprintf(fileListURLTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.ProjectID, personaID, withProgressStatus)
	return f.callFileRequestAndHandle(fileListURL)
}

func (f *FileService) WorkspaceFilesWithLimit(personaID int, withProgressStatus bool, limit int) (*types.FileSearchResponse, error) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "WorkspaceFilesWithLimit "+strconv.Itoa(personaID))
	}()
	base := f.Config.GetAPIBase()
	fileListURL := fmt.Sprintf(fileListURLTemplateWithLimit, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.ProjectID, personaID, withProgressStatus, limit)
	return f.callFileRequestAndHandle(fileListURL)
}

func (f *FileService) callFileRequestAndHandle(getUserFiles string) (*types.FileSearchResponse, error) {
	fileBytesResponse, err := f.QordobaClient.GetFromServer(getUserFiles)
	if err != nil {
		return nil, err
	}
	var response types.FileSearchResponse
	err = json.Unmarshal(fileBytesResponse, &response)
	if err != nil {
		log.Errorf("error occurred on server response unmarshalling: %v", err)
		return nil, err
	}
	return &response, nil
}

// DownloadFile function retrieves file in workspace
func (f *FileService) DownloadFile(personaID int, fileName string, file *types.File) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "DownloadFile")
	}()
	base := f.Config.GetAPIBase()
	getFileContentURL := fmt.Sprintf(fileDownloadTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.ProjectID, personaID, file.FileID)
	f.handleDownloadedFile(getFileContentURL, fileName)
}

func (f *FileService) handleDownloadedFile(fileRemoteURL, fileName string) {
	fileBytesResponse, err := f.QordobaClient.GetFromServer(fileRemoteURL)
	if err != nil {
		log.Errorf("error occurred on file s download (url = %s)\n%v", fileRemoteURL, fileName, err)
		return
	}
	log.Infof("file '%v' was downloaded", fileName)
	f.Local.Write(fileName, fileBytesResponse)
}

// DownloadSourceFile function retrieves all source files in workspace
func (f *FileService) DownloadSourceFile(fileName string, file *types.File, withUpdates bool) {
	base := f.Config.GetAPIBase()
	getFileContentURL := fmt.Sprintf(sourceFileDownloadTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.ProjectID, file.FileID, withUpdates)
	f.handleDownloadedFile(getFileContentURL, fileName)
}

// FindFile function
func (f *FileService) FindFile(fileName, version string) (*types.File, int) {
	log.Debugf("FindFile was called for file '%v'('%v')", fileName, version)
	workspace, err := f.WorkspaceService.LoadWorkspace()
	if err != nil {
		return nil, 0
	}
	base := f.Config.GetAPIBase()
	for _, persona := range workspace.Workspace.TargetPersonas {
		fileListURL := fmt.Sprintf(fileSearchURLTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.ProjectID, persona.ID, false, fileName, version)
		fileSearchResponse, err := f.callFileRequestAndHandle(fileListURL)
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

// deleteFoundFile function retrieve file and delete it remotedly
func (f *FileService) DeleteFile(fileName, version string) {
	log.Debugf("deleteFoundFile was called for file '%v'('%v')", fileName, version)
	file, _ := f.FindFile(fileName, version)
	if file != nil {
		f.deleteFoundFile(file)
	}
}

// deleteFoundFile func delete file from parameters
func (f *FileService) deleteFoundFile(file *types.File) {
	base := f.Config.GetAPIBase()
	deleteFileURL := fmt.Sprintf(fileDeleteTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.ProjectID, file.FileID)
	bytes, err := f.QordobaClient.DeleteFromServer(deleteFileURL)
	if err != nil {
		return
	}
	var deleteResponse types.FileDeleteResponse
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
