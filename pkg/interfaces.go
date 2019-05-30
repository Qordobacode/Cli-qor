package pkg

import (
	"github.com/qordobacode/cli-v2/pkg/types"
	"net/http"
)

type QordobaClient interface {
	GetFromServer(getURL string) ([]byte, error)
	PostToServer(postURL string, requestBody interface{}) (*http.Response, error)
	PutToServer(postURL string, requestBody interface{}) (*http.Response, error)
	DeleteFromServer(deleteURL string) ([]byte, error)
}

type Local interface {
	Read(path string) ([]byte, error)
	Write(fileName string, fileBytesResponse []byte)
	BuildFileName(file *types.File, suffix string) string
	FileExists(path string) bool
	QordobaHome() (string, error)
	PutInHome(fileName string, body []byte)
	LoadCached(cachedFileName string) ([]byte, error)
	FilesInFolder(folderPath string) []string
}

type ConfigurationService interface {
	ReadConfigInPath(path string) (*types.Config, error)
	LoadConfig() (*types.Config, error)
	SaveMainConfig(config *types.Config)
}

type WorkspaceService interface {
	LoadWorkspace() (*types.WorkspaceData, error)
}

type FileService interface {
	WorkspaceFiles(personaID int, withProgressStatus bool) (*types.FileSearchResponse, error)
	WorkspaceFilesWithLimit(personaID int, withProgressStatus bool, limit int) (*types.FileSearchResponse, error)
	FindFile(fileName, version string, withProgressStatus bool) (*types.File, int)
	DownloadFile(personaID int, fileName string, file *types.File)
	DownloadSourceFile(fileName string, file *types.File, withUpdates bool)
	PushFolder(folder, version string)
	PushFiles(fileList []string, version string)
	DeleteFile(fileName, version string)
}

type SegmentService interface {
	AddKey(fileName, version string, keyAddRequest *types.KeyAddRequest)
	UpdateKey(fileName, version string, keyAddRequest *types.KeyAddRequest)
	FindSegment(base, segmentName string, personaID int, file *types.File) *types.Segment
}
