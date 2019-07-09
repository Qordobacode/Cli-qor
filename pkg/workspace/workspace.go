package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"strconv"
	"time"
)

const (
	getWorkspacesTemnplate = "%s/v3/organizations/%d/workspaces"
	workspaceFileName      = "workspace.json"
)

// Service implements pkg.Service
type Service struct {
	Config        *types.Config
	QordobaClient pkg.QordobaClient
	Local         pkg.Local
}

// LoadWorkspace function retrieves a workspace
func (w *Service) LoadWorkspace() (*types.WorkspaceData, error) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "LoadWorkspace "+strconv.Itoa(int(w.Config.Qordoba.WorkspaceID)))
	}()
	workspaceResponse, err := w.cachedWorkspace()
	if err != nil || workspaceResponse == nil {
		workspaceResponse, err = w.workspacesFromServer()
	}
	if err != nil {
		return nil, err
	}
	for _, workspaceData := range workspaceResponse.Workspaces {
		if workspaceData.Workspace.ID == int(w.Config.Qordoba.WorkspaceID) {
			return &workspaceData, nil
		}
	}
	return nil, errors.New("workspace with id=" + string(w.Config.Qordoba.WorkspaceID) + " was not found")
}

// cachedWorkspace function returns cached workspace if it present AND still valid (invalidation period for
// cache is `invalidationPeriod`
func (w *Service) cachedWorkspace() (*types.WorkspaceResponse, error) {
	bodyBytes, err := w.Local.LoadCached(workspaceFileName)
	if err != nil {
		return nil, err
	}
	var workspaceResponse types.WorkspaceResponse
	err = json.Unmarshal(bodyBytes, &workspaceResponse)
	if err != nil {
		log.Errorf("error occurred on cached workspace read: %v", err)
		return nil, err
	}
	return &workspaceResponse, nil
}

// workspacesFromServer function retrieve list of all workspaces
func (w *Service) workspacesFromServer() (*types.WorkspaceResponse, error) {
	base := w.Config.GetAPIBase()
	// retrieve from server list of workspaces
	workspaceRequestURL := fmt.Sprintf(getWorkspacesTemnplate, base, w.Config.Qordoba.OrganizationID)
	bodyBytes, err := w.QordobaClient.GetFromServer(workspaceRequestURL)
	if err != nil {
		return nil, err
	}
	var workspaceResponse types.WorkspaceResponse
	err = json.Unmarshal(bodyBytes, &workspaceResponse)
	if err != nil {
		log.Errorf("error occurred on request for workspace: %v", err)
		return nil, err
	}
	w.Local.PutInHome(workspaceFileName, bodyBytes)
	return &workspaceResponse, nil
}
