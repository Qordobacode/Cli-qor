package file

import (
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestService_PushFolder(t *testing.T) {
	filesList := []string{"test.yaml"}
	service := buildFileService(t)
	local.EXPECT().FilesInFolder(gomock.Any()).Return(filesList)
	service.PushFolder(".", "")
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
