package pkg

import (
	"github.com/qordobacode/cli-v2/pkg/types"
	"net/http"
)

// QordobaClient interface collect all web-request related logic
type QordobaClient interface {
	GetFromServer(getURL string) ([]byte, error)
	PostToServer(postURL string, requestBody interface{}) (*http.Response, error)
	PutToServer(putURL string, requestBody interface{}) (*http.Response, error)
	DeleteFromServer(deleteURL string) ([]byte, error)
}

// Local interface collect all os and stdin-related logic
type Local interface {
	Read(path string) ([]byte, error)
	Write(fileName string, fileBytesResponse []byte)
	BuildDirectoryFilePath(j *types.File2Download, patterns []string, suffix string) string
	FileExists(path string) bool
	QordobaHome() (string, error)
	PutInHome(fileName string, body []byte)
	LoadCached(cachedFileName string) ([]byte, error)
	FilesInFolder(folderPath string) []string
	RenderTable2Stdin(header []string, data [][]string)
}

// ConfigurationService contains all methods about app configuration
type ConfigurationService interface {
	ReadConfigInPath(path string) (*types.Config, error)
	LoadConfig() (*types.Config, error)
	SaveMainConfig(config *types.Config)
}

// WorkspaceService contain workspace-related functionality
type WorkspaceService interface {
	LoadWorkspace() (*types.WorkspaceData, error)
}

// FileService contains all logic related to Qordoba's file
type FileService interface {
	WorkspaceFiles(personaID int, withProgressStatus bool) (*types.FileSearchResponse, error)
	WorkspaceFilesWithLimit(personaID int, withProgressStatus bool, limit int) (*types.FileSearchResponse, error)
	FindFile(fileName, version string, withProgressStatus bool) (*types.File, int)
	DownloadFile(personaID int, fileName string, file *types.File)
	DownloadSourceFile(fileName string, file *types.File, withUpdates bool)
	PushFolder(folder, version string)
	PushFiles(fileList []string, version string)
	DeleteFile(fileName, version string)
	FileScore(filename, version string) *types.ScoreResponseBody
}

// SegmentService contains all logic about Qordoba's segments
type SegmentService interface {
	FindSegment(fileName, fileVersion, key string) (*types.Segment, *types.File)
	AddKey(fileName, version string, keyAddRequest *types.KeyAddRequest)
	UpdateKey(fileName, version string, keyAddRequest *types.KeyAddRequest)
	DeleteKey(fileName, version, segmentKey string)
}
