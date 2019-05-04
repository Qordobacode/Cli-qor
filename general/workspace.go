package general

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qordobacode/cli-v2/log"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const (
	getWorkspacesTemnplate = "%s/v3/organizations/%d/workspaces"
	workspaceName  = "workspace.json"
)

// GetWorkspace function retrieves a workspace
func GetWorkspace(qordobaConfig *Config) (*Workspace, error) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "GetWorkspace " + strconv.Itoa(int(qordobaConfig.Qordoba.ProjectID)))
	}()
	workspaceResponse, err := GetCachedWorkspace()
	if err != nil || workspaceResponse == nil{
		workspaceResponse, err = GetAllWorkspaces(qordobaConfig)
	}
	if err != nil {
		return nil, err
	}
	for _, workspaceData := range workspaceResponse.Workspaces {
		if workspaceData.Workspace.ID == int(qordobaConfig.Qordoba.ProjectID) {
			return &workspaceData.Workspace, nil
		}
	}
	return nil, errors.New("workspace with id=" + string(qordobaConfig.Qordoba.ProjectID) + " was not found")
}

func GetCachedWorkspace() (*WorkspaceResponse, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("error occurred on home dir retrieval: %v", err)
		return nil, err
	}
	workspaceFilePath := getQordobaHomeDir(home) + string(os.PathSeparator) + workspaceName
	file, err := os.Stat(workspaceFilePath)
	if err != nil {
		return nil, err
	}
	modifiedtime := file.ModTime()
	// don't use cached workspace if 1 day has came
	if modifiedtime.Add(time.Hour * 24).Before(time.Now()) {
		return nil, errors.New("outdated file")
	}

	if !FileExists(workspaceFilePath) {
		log.Debugf("workspace not found: %v", workspaceFilePath)
		return nil, fmt.Errorf("cached workspace file was not found")
	}
	// read config from file
	bodyBytes, err := ioutil.ReadFile(workspaceFilePath)
	if err != nil {
		log.Debugf("file not found: %v", err)
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
