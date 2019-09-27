package file

import (
	"github.com/golang/mock/gomock"
	"github.com/qordobacode/cli-v2/pkg/types"
	"testing"
)

func startDownload(t *testing.T) *Service {
	service := buildFileService(t)
	local.EXPECT().Write("testing.yaml", gomock.Any())
	return service
}

func Test_DownloadFile(t *testing.T) {
	service := startDownload(t)
	file := types.File{}
	person := types.Person{
		Code:      "en-us",
		Direction: "",
		ID:        100,
		Name:      "",
	}
	service.DownloadFile(person, "testing.yaml", &file)

}
