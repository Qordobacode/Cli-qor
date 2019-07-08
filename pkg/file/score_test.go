package file

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func startScore(t *testing.T) *Service {
	resp := `{
  "snapshotTime": 0,
  "documentScore": 0,
  "breakdown": [
    {
      "category": "test",
      "issueCount": 0,
      "score": 0,
      "enabled": true
    }
  ]
}`
	service := buildFileService(t)
	filesList := []string{"test.yaml"}
	local.EXPECT().FilesInFolder(gomock.Any()).Return(filesList)
	client.EXPECT().GetFromServer("https://app.qordoba.com/v3/contentscore/organizations/0/workspaces/0/documents/7637/personas/100/score?documentLength=1").
		Return([]byte(resp), nil)
	return service
}

func TestService_FileScore(t *testing.T) {
	service := startScore(t)
	service.FileScore("test.json", "")
}
