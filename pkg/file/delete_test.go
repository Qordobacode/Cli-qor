package file

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func prepareDeleteTest(t *testing.T) *Service {
	service := buildFileService(t)
	resp := `{
  "success": true
}`
	client.EXPECT().DeleteFromServer(gomock.Any()).Return([]byte(resp), nil)
	return service
}

func TestService_DeleteFileNotFound(t *testing.T) {
	service := prepareDeleteTest(t)
	service.DeleteFile("test.yaml", "")
}

func TestService_DeleteFile(t *testing.T) {
	service := prepareDeleteTest(t)
	service.DeleteFile("test.json", "")
}
