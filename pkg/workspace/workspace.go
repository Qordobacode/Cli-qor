package workspace

import (
	"errors"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"time"
)

const (
	getWorkspacesTemplate = "%s/v3/organizations/%d/workspaces?limit=%d&offset=%d"
	limit                 = 500
	workspaceFileName     = "workspace.json"
)

var (
	workspaceCacheWasUpdated = false
)

// Service implements pkg.Service
type Service struct {
	Config        *types.Config
	QordobaClient pkg.QordobaClient
	Local         pkg.Local
}

// LoadWorkspace function retrieves a workspace
func (w *Service) LoadWorkspace() (*types.WorkspaceData, error) {
	workspaceResponse, err := w.cachedWorkspace()
	if err == nil && workspaceResponse != nil {
		for _, workspaceData := range workspaceResponse.Workspaces {
			if workspaceData.Workspace.ID == int(w.Config.Qordoba.WorkspaceID) {
				return &workspaceData, nil
			}
		}
	}
	return w.WorkspaceFromServer()
}

func (w *Service) WorkspaceFromServer() (*types.WorkspaceData, error) {
	if workspaceCacheWasUpdated {
		return nil, errors.New("workspace has been already updated")
	}
	workspaceResponse, err := w.loadServerWorkspaceResponse()
	if err == nil && workspaceResponse != nil {
		for _, workspaceData := range workspaceResponse.Workspaces {
			if workspaceData.Workspace.ID == int(w.Config.Qordoba.WorkspaceID) {
				return &workspaceData, nil
			}
		}
	} else {
		return nil, err
	}
	err = fmt.Errorf("workspace with id=%v was not found", w.Config.Qordoba.WorkspaceID)
	log.Errorf(err.Error())
	return nil, err
}

// cachedWorkspace function returns cached workspace if it present AND still valid (invalidation period for
// cache is `invalidationPeriod`
func (w *Service) cachedWorkspace() (*types.WorkspaceResponse, error) {
	bodyBytes, err := w.Local.LoadCached(workspaceFileName)
	if err != nil {
		return nil, err
	}
	var workspaceResponse types.WorkspaceResponse
	err = workspaceResponse.UnmarshalJSON(bodyBytes)
	if err != nil {
		log.Errorf("error occurred on cached workspace read: %v", err)
		return nil, err
	}
	return &workspaceResponse, nil
}

// loadServerWorkspaceResponse function retrieve list of all workspaces
func (w *Service) loadServerWorkspaceResponse() (*types.WorkspaceResponse, error) {
	log.Infof("start to download organization's workspace structure...")
	start := time.Now()
	base := w.Config.GetAPIBase()
	result := &types.WorkspaceResponse{
		Meta: types.Meta{
			Paging: types.Paging{
				TotalEnabled: 0,
				TotalResults: 1,
			},
		},
		Workspaces: make([]types.WorkspaceData, 0),
	}

	errNum := 0
	for offset := 0; offset < result.Meta.Paging.TotalResults; offset += limit {
		// retrieve from server list of workspaces
		workspaceRequestURL := fmt.Sprintf(getWorkspacesTemplate, base, w.Config.Qordoba.OrganizationID, limit, offset)
		bodyBytes, err := w.QordobaClient.GetFromServer(workspaceRequestURL)
		if err != nil {
			if errNum == 0 {
				// try to repeat failed request 1 time
				offset -= limit
				errNum++
			}
			continue
		}
		var workspaceResponse types.WorkspaceResponse
		err = workspaceResponse.UnmarshalJSON(bodyBytes)
		if err != nil {
			log.Errorf("error occurred on request for workspace: %v", err)
			continue
		}
		result.Meta.Paging.TotalResults = workspaceResponse.Meta.Paging.TotalResults
		result.Meta.Paging.TotalEnabled = workspaceResponse.Meta.Paging.TotalEnabled
		result.Workspaces = append(result.Workspaces, workspaceResponse.Workspaces...)
		errNum = 0
		elapsed := time.Since(start)
		log.Infof("%v. Downloaded %d/%d organization's workspaces", elapsed, len(result.Workspaces), result.Meta.Paging.TotalResults)
	}

	bytes, err := result.MarshalJSON()
	if err == nil {
		w.Local.PutInHome(workspaceFileName, bytes)
	}
	workspaceCacheWasUpdated = true
	return result, nil
}
