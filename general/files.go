package general

import "fmt"

const (
	fileListURLTemplate = "%s/v3/organizations/%d/workspaces/%d/personas/%d/files"
)

// GetFilesInWorkspace function retrieves all files in workspace
func GetFilesInWorkspace(config *Config, personaID int) {
	base := config.GetAPIBase()
	getUserFiles := fmt.Sprintf(fileListURLTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, personaID)
	response, err := GetFromServer(config, getUserFiles)
}
