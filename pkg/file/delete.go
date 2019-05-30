package file

import (
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
)

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