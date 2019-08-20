package file

import (
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"time"
)

// DownloadFile function retrieves file in workspace
func (f *Service) DownloadFile(personaID int, fileName string, file *types.File) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "DownloadFile")
	}()
	base := f.Config.GetAPIBase()
	getFileContentURL := fmt.Sprintf(fileDownloadTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.WorkspaceID, personaID, file.FileID)
	f.handleDownloadedFile(getFileContentURL, fileName)
}

func (f *Service) handleDownloadedFile(fileRemoteURL, fileName string) {
	fileBytesResponse, err := f.QordobaClient.GetFromServer(fileRemoteURL)
	if err != nil {
		log.Errorf("error occurred on file %s download (url = %s)\n%v", fileName, fileRemoteURL, err.Error())
		return
	}
	log.Infof("file '%s' was downloaded", fileName)
	f.Local.Write(fileName, fileBytesResponse)
}

// DownloadSourceFile function retrieves all source files in workspace
func (f *Service) DownloadSourceFile(fileName string, file *types.File, withUpdates bool) {
	base := f.Config.GetAPIBase()
	getFileContentURL := fmt.Sprintf(sourceFileDownloadTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.WorkspaceID, file.FileID, withUpdates)
	f.handleDownloadedFile(getFileContentURL, fileName)
}
