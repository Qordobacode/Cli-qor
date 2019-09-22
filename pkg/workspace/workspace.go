package workspace

import (
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
	workspaceResponse, err = w.workspacesFromServer()
	if err == nil && workspaceResponse != nil {
		for _, workspaceData := range workspaceResponse.Workspaces {
			if workspaceData.Workspace.ID == int(w.Config.Qordoba.WorkspaceID) {
				return &workspaceData, nil
			}
		}
	}
	err = fmt.Errorf("workspace with id=%v was not found", string(w.Config.Qordoba.WorkspaceID))
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

// workspacesFromServer function retrieve list of all workspaces
func (w *Service) workspacesFromServer() (*types.WorkspaceResponse, error) {
	log.Infof("start to download %v organization's workspace structure...", w.Config.Qordoba.WorkspaceID)
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
		log.Infof("request %d/%d...", offset/limit+1, (result.Meta.Paging.TotalResults+limit-1)/limit)
	}

	bytes, err := result.MarshalJSON()
	if err == nil {
		w.Local.PutInHome(workspaceFileName, bytes)
	}
	elapsed := time.Since(start)
	log.Infof("Organization %d workspace structure is downloaded in %v. %d workspaces were downloaded", w.Config.Qordoba.WorkspaceID, elapsed, len(result.Workspaces))
	return result, nil
}
