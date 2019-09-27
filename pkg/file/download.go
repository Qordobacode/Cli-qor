package file

import (
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"os"
	"path/filepath"
	"time"
)

// DownloadFile function retrieves file in workspace
func (f *Service) DownloadFile(persona types.Person, fileName string, file *types.File) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "DownloadFile")
	}()
	base := f.Config.GetAPIBase()
	getFileContentURL := fmt.Sprintf(fileDownloadTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.WorkspaceID, persona.ID, file.FileID)
	f.handleDownloadedFile(getFileContentURL, fileName, persona.Code)
}

func (f *Service) handleDownloadedFile(fileRemoteURL, fileName, language string) {
	fileBytesResponse, err := f.QordobaClient.GetFromServer(fileRemoteURL)
	if err != nil {
		log.Errorf("error occurred on file %s download (url = %s)\n%v", fileName, fileRemoteURL, err.Error())
		return
	}
	if len(f.Config.Push.Sources.Folders) > 0 {
		fileName = filepath.Join(f.Config.Push.Sources.Folders[0], fileName)
	}
	err = os.MkdirAll(filepath.Dir(fileName), 0755)
	if err != nil {
		log.Errorf("error occurred on creating new directories")
	}
	f.Local.Write(fileName, fileBytesResponse)
	if language == "" {
		log.Infof("file %s was downloaded", fileName)
	} else {
		log.Infof("file %s was downloaded for language %s", fileName, language)
	}
}

// DownloadSourceFile function retrieves all source files in workspace
func (f *Service) DownloadSourceFile(fileName string, file *types.File, withUpdates bool) {
	base := f.Config.GetAPIBase()
	getFileContentURL := fmt.Sprintf(sourceFileDownloadTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.WorkspaceID, file.FileID, withUpdates)
	f.handleDownloadedFile(getFileContentURL, fileName, "")
}
