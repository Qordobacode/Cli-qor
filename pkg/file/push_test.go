package file

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestService_PushFolder(t *testing.T) {
	filesList := []string{"test.yaml"}
	service := buildFileService(t)
	local.EXPECT().FilesInFolder(gomock.Any()).Return(filesList)
	service.PushFolder(".", "")
}
