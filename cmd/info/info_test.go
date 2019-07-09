package info

import (
	"github.com/golang/mock/gomock"
	"github.com/qordobacode/cli-v2/pkg/file"
	"github.com/qordobacode/cli-v2/pkg/mock"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/qordobacode/cli-v2/pkg/workspace"
	"testing"
)

var (
	clientMock        *mock.MockQordobaClient
	workspaceResponse = `{"meta":{"paging":{"totalResults":14}},"workspaces":[{"workspace":{"contentTypeCodes":[{"extensions":["json"],"name":"JSON"}],"createdBy":{"id":6,"name":"May Habib","role":""},"createdOn":1560424936000,"id":365,"name":"DemoTesting","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"}],"timezone":"PST8PDT"},"workflow":[{"id":682,"name":"Under Translation","order":0,"complete":false},{"id":683,"name":"Key Management","order":0,"complete":false},{"id":684,"name":"Under Edit","order":1,"complete":false},{"id":685,"name":"Review","order":1,"complete":false},{"id":686,"name":"Review","order":2,"complete":false},{"id":687,"name":"Staging","order":3,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["txt"],"name":"txt all text"},{"extensions":["json"],"name":"JSON"},{"extensions":["xml"],"name":"Android XML"},{"extensions":["yml","yaml"],"name":"YAML with MD"}],"createdBy":{"id":1103,"name":"Evgenii Morenkov","role":""},"createdOn":1558178398000,"id":333,"name":"emorenkov-test","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"fr-fr","direction":"ltr","id":110,"name":"French - France"},{"code":"es-es","direction":"ltr","id":234,"name":"Spanish - Spain"},{"code":"de-de","direction":"ltr","id":124,"name":"German - Germany"}],"timezone":"PST8PDT"},"workflow":[{"id":620,"name":"Translation","order":0,"complete":false},{"id":621,"name":"Review","order":1,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["txt"],"name":"txt all text"},{"extensions":["md","text"],"name":"Markdown"}],"createdBy":{"id":6,"name":"May Habib","role":""},"createdOn":1553122998000,"id":288,"name":"jrMaxHistory","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"en-gb","direction":"ltr","id":92,"name":"English - United Kingdom"}],"timezone":"PST8PDT"},"workflow":[{"id":535,"name":"Content Review","order":0,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["json"],"name":"JSON"}],"createdBy":{"id":6,"name":"May Habib","role":""},"createdOn":1559205338000,"id":348,"name":"husam-otri/jspage","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"ar-dz","direction":"rtl","id":3,"name":"Arabic - Algeria"}],"timezone":"PST8PDT"},"workflow":[{"id":650,"name":"Translation","order":0,"complete":false},{"id":651,"name":"Review","order":1,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["txt"],"name":"txt all text"},{"extensions":["json"],"name":"JSON"},{"extensions":["yml","yaml"],"name":"YAML with MD"},{"extensions":["xml"],"name":"Android XML"}],"createdBy":{"id":1103,"name":"Evgenii Morenkov","role":""},"createdOn":1558178413000,"id":334,"name":"emorenkov-test","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"de-de","direction":"ltr","id":124,"name":"German - Germany"},{"code":"es-es","direction":"ltr","id":234,"name":"Spanish - Spain"},{"code":"fr-fr","direction":"ltr","id":110,"name":"French - France"}],"timezone":"PST8PDT"},"workflow":[{"id":622,"name":"Translation","order":0,"complete":false},{"id":623,"name":"Review","order":1,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["skt"],"name":"Sketch"}],"createdBy":{"id":6,"name":"May Habib","role":""},"createdOn":1542897855000,"id":188,"name":"test eloqua","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-au","direction":"ltr","id":72,"name":"English - Australia"},"targetPersonas":[{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"}],"timezone":"PST8PDT"},"workflow":[{"description":"First Milestone Description","id":402,"name":"ONE","order":0,"complete":false},{"description":"New Milestone Description","id":403,"name":"TWO","order":1,"complete":false},{"description":"New Milestone Description","id":404,"name":"THREE","order":2,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["json"],"name":"JSON"}],"createdBy":{"id":6,"name":"May Habib","role":""},"createdOn":1559043462000,"id":339,"name":"husam-otri/allFilesTypes","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"es-pe","direction":"ltr","id":230,"name":"Spanish - Peru"}],"timezone":"PST8PDT"},"workflow":[{"id":632,"name":"Translation","order":0,"complete":false},{"id":633,"name":"Review","order":1,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["json"],"name":"JSON"}],"createdBy":{"id":1022,"name":"Hussam Otri","role":""},"createdOn":1537797829000,"id":148,"name":"test wf","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":315,"name":"Inglize"},"targetPersonas":[{"code":"es-pe","direction":"ltr","id":230,"name":"Spanish - Peru"}],"timezone":"PST8PDT"},"workflow":[{"description":"First Milestone Description","id":335,"name":"NOT ONE","order":0,"complete":false},{"description":"New Milestone Description","id":336,"name":"TWO","order":1,"complete":false},{"description":"New Milestone Description","id":337,"name":"THREE","order":2,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["txt"],"name":"txt all text"},{"extensions":["xml"],"name":"Android XML"},{"extensions":["json"],"name":"JSON"},{"extensions":["yml","yaml"],"name":"YAML with MD"}],"createdBy":{"id":1103,"name":"Evgenii Morenkov","role":""},"createdOn":1558178506000,"id":335,"name":"emorenkov-test","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"fr-fr","direction":"ltr","id":110,"name":"French - France"},{"code":"de-de","direction":"ltr","id":124,"name":"German - Germany"},{"code":"es-es","direction":"ltr","id":234,"name":"Spanish - Spain"}],"timezone":"PST8PDT"},"workflow":[{"id":624,"name":"Translation","order":0,"complete":false},{"id":625,"name":"Review","order":1,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["json"],"name":"JSON"}],"createdBy":{"id":6,"name":"May Habib","role":""},"createdOn":1543359231000,"id":194,"name":"test hussam","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"en-ca","direction":"ltr","id":76,"name":"English - Canada"}],"timezone":"PST8PDT"},"workflow":[{"id":412,"name":"Translation","order":0,"complete":false},{"id":413,"name":"Review","order":1,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["txt"],"name":"txt all text"},{"extensions":["yml","yaml"],"name":"YAML with MD"},{"extensions":["json"],"name":"JSON"},{"extensions":["xml"],"name":"Android XML"}],"createdBy":{"id":1103,"name":"Evgenii Morenkov","role":""},"createdOn":1558178519000,"id":336,"name":"emorenkov-test","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"fr-fr","direction":"ltr","id":110,"name":"French - France"},{"code":"de-de","direction":"ltr","id":124,"name":"German - Germany"},{"code":"es-es","direction":"ltr","id":234,"name":"Spanish - Spain"}],"timezone":"PST8PDT"},"workflow":[{"id":626,"name":"Translation","order":0,"complete":false},{"id":627,"name":"Review","order":1,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["json"],"name":"JSON"}],"createdBy":{"id":6,"name":"May Habib","role":""},"createdOn":1532604908000,"id":99,"name":"File Complete Status","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"pl-pl","direction":"ltr","id":180,"name":"Polish - Poland"},{"code":"es-es","direction":"ltr","id":234,"name":"Spanish - Spain"},{"code":"pt-pt","direction":"ltr","id":184,"name":"Portuguese - Portugal"},{"code":"en-gb","direction":"ltr","id":92,"name":"English - United Kingdom"},{"code":"fr-fr","direction":"ltr","id":110,"name":"French - France"}],"timezone":"PST8PDT"},"workflow":[{"description":"First Milestone Description","id":247,"name":"ONE","order":0,"complete":false},{"description":"New Milestone Description","id":248,"name":"TWO","order":1,"complete":false},{"description":"New Milestone Description","id":249,"name":"THREE","order":2,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["txt"],"name":"txt all text"},{"extensions":["xml"],"name":"Android XML"},{"extensions":["json"],"name":"JSON"},{"extensions":["yml","yaml"],"name":"YAML with MD"}],"createdBy":{"id":1103,"name":"Evgenii Morenkov","role":""},"createdOn":1558178353000,"id":332,"name":"emorenkov-test","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"de-de","direction":"ltr","id":124,"name":"German - Germany"},{"code":"es-es","direction":"ltr","id":234,"name":"Spanish - Spain"},{"code":"fr-fr","direction":"ltr","id":110,"name":"French - France"}],"timezone":"PST8PDT"},"workflow":[{"id":618,"name":"Translation","order":0,"complete":false},{"id":619,"name":"Review","order":1,"complete":false}]},{"workspace":{"contentTypeCodes":[{"extensions":["json"],"name":"JSON"}],"createdBy":{"id":6,"name":"May Habib","role":""},"createdOn":1560445696000,"id":369,"name":"husam-otri/anotherTest","organizationId":9,"segmentation":"default","sourcePersona":{"code":"en-us","direction":"ltr","id":94,"name":"English - United States"},"targetPersonas":[{"code":"de-de","direction":"ltr","id":124,"name":"German - Germany"}],"timezone":"PST8PDT"},"workflow":[{"id":694,"name":"Translation","order":0,"complete":false},{"id":695,"name":"Review","order":1,"complete":false}]}]}`
)

func startConfig(t *testing.T) {
	controller := gomock.NewController(t)
	workspaceMock := mock.NewMockWorkspaceService(controller)
	local := mock.NewMockLocal(controller)
	clientMock = mock.NewMockQordobaClient(controller)
	appConfig = &types.Config{
		Qordoba: types.QordobaConfig{
			AudienceMap: map[string]string{
				"pl-pl": "test",
			},
		},
		BaseURL: "baseURL",
	}
	workspaceData := &types.WorkspaceData{
		Workspace: types.Workspace{
			TargetPersonas: []types.Person{
				{ID: 365},
			},
		},
	}
	workspaceMock.EXPECT().LoadWorkspace().Return(workspaceData, nil)
	clientMock.EXPECT().GetFromServer(gomock.Any()).Return([]byte(`getResponse`), nil)
	clientMock.EXPECT().DeleteFromServer(gomock.Any()).Return([]byte(`deleteResponse`), nil)
	fileService = &file.Service{
		Config:           &types.Config{},
		WorkspaceService: workspaceMock,
		Local:            local,
		QordobaClient:    clientMock,
	}
	local.EXPECT().LoadCached(gomock.Any()).Return([]byte(workspaceResponse), nil)
	workspaceService = &workspace.Service{
		Config: &types.Config{
			Qordoba: types.QordobaConfig{
				WorkspaceID: 365,
			},
		},
		QordobaClient: clientMock,
		Local:         local,
	}
}

func Test_startLocalServices(t *testing.T) {
	startLocalServices(nil, nil)
}

func TestNewLsCommand(t *testing.T) {
	lsCommand := NewLsCommand()
	startConfig(t)
	printLs(lsCommand, []string{})
}

func TestNewScoreCommand(t *testing.T) {
	scoreCommand := NewScoreCommand()
	startConfig(t)
	scoreFile(scoreCommand, []string{"result.yaml"})
}

func TestNewStatusCommand(t *testing.T) {
	statusCommand := NewStatusCommand()
	startConfig(t)
	runStatus(statusCommand, []string{"result.yaml"})
}
