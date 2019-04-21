package general

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qordobacode/cli-v2/log"
	"io/ioutil"
	"net/http"
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
	response, err := GetFromServer(qordobaConfig, workspaceRequestURL)
	if err != nil {
		log.Errorf("error occurred on workspace get request: %v", err)
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("error occurred on body read: %v", err)
		return nil, err
	}
	if response.StatusCode/100 != 2 {
		if response.StatusCode == http.StatusUnauthorized {
			log.Errorf("User is not authorised for this request. Check `access_token` in configuration.")
		} else {
			log.Errorf("Error occurred on workspace request. Status: %d, Response : %v", response.Status, string(bodyBytes))
		}
		return nil, errors.New("unsuccessful request")
	}
	var workspaceResponse WorkspaceResponse
	err = json.Unmarshal(bodyBytes, &workspaceResponse)
	if err != nil {
		log.Errorf("error occurred on request for workspace: %v", err)
		return nil, err
	}
	return &workspaceResponse, nil
}
