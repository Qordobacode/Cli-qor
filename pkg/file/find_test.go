package file

import (
	"github.com/golang/mock/gomock"
	"github.com/qordobacode/cli-v2/pkg/mock"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var (
	appConfig        *types.Config
	workspaceService *mock.MockWorkspaceService
	local            *mock.MockLocal
	client           *mock.MockQordobaClient
)

func buildFileService(t *testing.T) *Service {
	appConfig = &types.Config{}
	controller := gomock.NewController(t)
	workspaceService = mock.NewMockWorkspaceService(controller)
	workspaceData := &types.WorkspaceData{
		Workspace: types.Workspace{
			TargetPersonas: []types.Person{
				{ID: 100},
			},
		},
	}
	workspaceService.EXPECT().LoadWorkspace().Return(workspaceData, nil)
	client = mock.NewMockQordobaClient(controller)
	local = mock.NewMockLocal(controller)

	file, err := ioutil.ReadFile("filesearch_response.json")
	client.EXPECT().GetFromServer(gomock.Any()).Return(file, err)

	fileService := &Service{
		Config:           appConfig,
		WorkspaceService: workspaceService,
		Local:            local,
		QordobaClient:    client,
	}
	return fileService
}

func TestService_FindFileNoName(t *testing.T) {
	service := buildFileService(t)
	file, personID := service.FindFile("", "", false)
	assert.Nil(t, file)
	assert.Equal(t, 0, personID)
}

func TestService_FindFile(t *testing.T) {
	service := buildFileService(t)
	file, personID := service.FindFile("test.json", "", false)
	assert.NotNil(t, file)
	assert.Equal(t, 100, personID)
}

func TestService_FindFileNotFound(t *testing.T) {
	service := buildFileService(t)
	file, personID := service.FindFile("test.json", "version-1", false)
	assert.Nil(t, file)
	assert.Equal(t, 0, personID)
}

func TestService_WorkspaceFiles(t *testing.T) {
	service := buildFileService(t)
	file, err := service.WorkspaceFiles(100, false)
	assert.NotNil(t, file)
	assert.Nil(t, err)
}

func TestService_WorkspaceFilesWithLimit(t *testing.T) {
	service := buildFileService(t)
	file, err := service.WorkspaceFilesWithLimit(100, false, 100)
	assert.NotNil(t, file)
	assert.Nil(t, err)
}
