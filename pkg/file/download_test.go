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
	service.DownloadFile(100, "testing.yaml", &file)

}
