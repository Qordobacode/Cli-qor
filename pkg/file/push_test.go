package file

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestService_PushFolder(t *testing.T) {
	filesList := []string{"test.yaml"}
	service := buildFileService(t)
	local.EXPECT().FilesInFolder(gomock.Any()).Return(filesList)
	service.PushFolder(".", "", false)
}

func TestService_PushFiles(t *testing.T) {
	filesList := []string{"filesearch_response.json", "notfound.json"}
	service := buildFileService(t)
	r := ioutil.NopCloser(strings.NewReader("some server response"))
	resp := http.Response{
		StatusCode: 200,
		Body:       r,
	}
	client.EXPECT().PostToServer(gomock.Any(), gomock.Any()).Return(&resp, nil)
	service.PushFiles(filesList, "v1")
}

func Test_Test(t *testing.T) {
	dir, _ := os.Getwd()
	relativeFilePath, _ := filepath.Rel(dir, `C:\data\code\Cli-qor\test\csv\core.csv`)
	fmt.Printf("path = %s", relativeFilePath)
}
