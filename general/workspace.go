package general

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qordobacode/cli-v2/log"
)

const (
	getWorkspacesTemnplate = "%s/v3/organizations/%d/workspaces"
)

// GetAllWorkspaces function retrieve list of all workspaces
func GetWorkspace(qordobaConfig *Config) (*Workspace, error) {
	allWorkspaces, err := GetAllWorkspaces(qordobaConfig)
	if err != nil {
		return nil, err
	}
	for _, workspaceData := range allWorkspaces.Workspaces {
		if workspaceData.Workspace.ID == int(qordobaConfig.Qordoba.ProjectID) {
			return &workspaceData.Workspace, nil
		}
	}

	return nil, errors.New("workspace with id=" + string(qordobaConfig.Qordoba.ProjectID) + " was not found")
}

// GetAllWorkspaces function retrieve list of all workspaces
func GetAllWorkspaces(qordobaConfig *Config) (*WorkspaceResponse, error) {
	base := qordobaConfig.GetAPIBase()
	// retrieve from server list of workspaces
	workspaceRequestURL := fmt.Sprintf(getWorkspacesTemnplate, base, qordobaConfig.Qordoba.OrganizationID)
	bodyBytes, err := GetFromServer(qordobaConfig, workspaceRequestURL)
	if err != nil {
		return nil, err
	}
	var workspaceResponse WorkspaceResponse
	err = json.Unmarshal(bodyBytes, &workspaceResponse)
	if err != nil {
		log.Errorf("error occurred on request for workspace: %v", err)
		return nil, err
	}
	return &workspaceResponse, nil
}
