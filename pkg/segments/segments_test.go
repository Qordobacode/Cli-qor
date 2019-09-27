package segments

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/qordobacode/cli-v2/pkg/mock"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var (
	fileService           *mock.MockFileService
	qordobaClient         *mock.MockQordobaClient
	workspaceService      *mock.MockWorkspaceService
	workspaceJSON         = `{ "workspace": { "contentTypeCodes": [ { "extensions": [ "json" ], "name": "JSON" } ], "createdBy": { "id": 6, "name": "May Habib", "role": "" }, "createdOn": 1532604908000, "id": 99, "name": "File Complete Status", "organizationId": 9, "segmentation": "default", "sourcePersona": { "code": "en-us", "direction": "ltr", "id": 94, "name": "English - United States" }, "targetPersonas": [ { "code": "pl-pl", "direction": "ltr", "id": 180, "name": "Polish - Poland" }, { "code": "es-es", "direction": "ltr", "id": 234, "name": "Spanish - Spain" }, { "code": "pt-pt", "direction": "ltr", "id": 184, "name": "Portuguese - Portugal" }, { "code": "en-gb", "direction": "ltr", "id": 92, "name": "English - United Kingdom" }, { "code": "fr-fr", "direction": "ltr", "id": 110, "name": "French - France" } ], "timezone": "PST8PDT" }, "workflow": [ { "description": "First Milestone Description", "id": 247, "name": "ONE", "order": 0, "complete": false }, { "description": "New Milestone Description", "id": 248, "name": "TWO", "order": 1, "complete": false }, { "description": "New Milestone Description", "id": 249, "name": "THREE", "order": 2, "complete": false } ] }`
	segmentSearchResponse = `{ "meta": { "paging": { "totalResults": 0 } }, "segments": [ { "lastSaved": 0, "order": 0, "pluralRule": "string", "plurals": "string", "reference": "some-reference", "segment": "string", "segmentId": 0, "ssMatch": 0, "ssText": "some-text", "stringKey": "/some-key", "target": "string", "targetId": 0 } ] }`
)

func startSegmentService(t *testing.T) *SegmentService {
	controller := gomock.NewController(t)
	fileService = mock.NewMockFileService(controller)
	qordobaClient = mock.NewMockQordobaClient(controller)
	workspaceService = mock.NewMockWorkspaceService(controller)
	file := types.File{
		Filename: "config.yaml",
		Version:  "v1",
	}
	fileService.EXPECT().FindFile("config.yaml", "v1", false).Return(&file, 100)
	fileService.EXPECT().FindFile("config.yaml", "v2", false).Return(nil, 0)
	var workspaceData types.WorkspaceData
	err := json.Unmarshal([]byte(workspaceJSON), &workspaceData)
	assert.Nil(t, err)
	workspaceService.EXPECT().LoadWorkspace().Return(&workspaceData, nil)
	qordobaClient.EXPECT().GetFromServer(gomock.Any()).Times(5).
		Return([]byte(segmentSearchResponse), nil)
	return &SegmentService{
		Config: &types.Config{
			Qordoba: types.QordobaConfig{},
		},
		QordobaClient:    qordobaClient,
		FileService:      fileService,
		WorkspaceService: workspaceService,
	}
}

func TestSegmentService_AddKeyFileNotFound(t *testing.T) {
	service := startSegmentService(t)
	keyAddRequest := &types.KeyAddRequest{
		Key:       "/key",
		Source:    "source",
		Reference: "reference",
	}
	service.AddKey("config.yaml", "v2", keyAddRequest)
}

func TestSegmentService_AddKey(t *testing.T) {
	service := startSegmentService(t)
	r := ioutil.NopCloser(strings.NewReader("some server response"))
	response := &http.Response{
		StatusCode: 200,
		Body:       r,
	}
	qordobaClient.EXPECT().PostToServer("https://app.qordoba.com/v3/organizations/0/workspaces/0/files/0/segments/keyAdd", gomock.Any()).
		Return(response, nil)
	keyAddRequest := &types.KeyAddRequest{
		Key:       "/key",
		Source:    "source",
		Reference: "reference",
	}
	service.AddKey("config.yaml", "v1", keyAddRequest)
}

func TestSegmentService_AddKeyBadResponse(t *testing.T) {
	service := startSegmentService(t)
	r := ioutil.NopCloser(strings.NewReader("some server response"))
	response := &http.Response{
		StatusCode: 400,
		Body:       r,
	}
	qordobaClient.EXPECT().PostToServer("https://app.qordoba.com/v3/organizations/0/workspaces/0/files/0/segments/keyAdd", gomock.Any()).
		Return(response, nil)
	keyAddRequest := &types.KeyAddRequest{
		Key:       "/key",
		Source:    "source",
		Reference: "reference",
	}
	service.AddKey("config.yaml", "v1", keyAddRequest)
}

func TestSegmentService_FindSegment(t *testing.T) {
	service := startSegmentService(t)
	segment, file := service.FindSegment("config.yaml", "v1", "/some-key")
	assert.NotNil(t, segment)
	assert.NotNil(t, file)
}

func TestSegmentService_FindSegmentNotFound(t *testing.T) {
	service := startSegmentService(t)
	segment, file := service.FindSegment("config.yaml", "v2", "/some-key")
	assert.Nil(t, segment)
	assert.Nil(t, file)
}

func TestSegmentService_UpdateKey(t *testing.T) {
	service := startSegmentService(t)
	keyAddRequest := &types.KeyAddRequest{
		Key:       "/some-key",
		Source:    "some-source",
		Reference: "some-reference",
	}
	r := ioutil.NopCloser(strings.NewReader("some server response"))
	response := &http.Response{
		StatusCode: 200,
		Body:       r,
	}
	qordobaClient.EXPECT().PutToServer(gomock.Any(), gomock.Any()).Times(5).
		Return(response, nil)
	service.UpdateKey("config.yaml", "v1", keyAddRequest)
}

func TestSegmentService_UpdateKeyErrorOnUpdate(t *testing.T) {
	service := startSegmentService(t)
	keyAddRequest := &types.KeyAddRequest{
		Key:       "/some-key",
		Source:    "some-source",
		Reference: "some-reference",
	}
	r := ioutil.NopCloser(strings.NewReader("some server response"))
	response := &http.Response{
		StatusCode: 400,
		Body:       r,
	}
	qordobaClient.EXPECT().PutToServer(gomock.Any(), gomock.Any()).Times(5).
		Return(response, nil)
	service.UpdateKey("config.yaml", "v1", keyAddRequest)
}

func TestSegmentService_DeleteKey(t *testing.T) {
	service := startSegmentService(t)
	qordobaClient.EXPECT().DeleteFromServer(gomock.Any()).
		Return([]byte("some-response"), nil)
	service.DeleteKey("config.yaml", "v1", "/some-key")
}

func TestSegmentService_DeleteKeyNotFoundFile(t *testing.T) {
	service := startSegmentService(t)
	qordobaClient.EXPECT().DeleteFromServer(gomock.Any()).
		Return([]byte("some-response"), nil)
	service.DeleteKey("config.yaml", "v2", "/some-key")
}

func Test_Test(t *testing.T) {
	service := startSegmentService(t)
	key := service.handleSegmentKey("/test")
	assert.Equal(t, key, "/test")
}
