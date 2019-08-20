package mock

import (
	"github.com/golang/mock/gomock"
	"github.com/qordobacode/cli-v2/pkg/types"
	"net/http"
	"reflect"
)

// MockQordobaClient is a mock of QordobaClient interface
type MockQordobaClient struct {
	ctrl     *gomock.Controller
	recorder *MockQordobaClientMockRecorder
}

// MockQordobaClientMockRecorder is the mock recorder for MockQordobaClient
type MockQordobaClientMockRecorder struct {
	mock *MockQordobaClient
}

// NewMockQordobaClient creates a new mock instance
func NewMockQordobaClient(ctrl *gomock.Controller) *MockQordobaClient {
	mock := &MockQordobaClient{ctrl: ctrl}
	mock.recorder = &MockQordobaClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockQordobaClient) EXPECT() *MockQordobaClientMockRecorder {
	return m.recorder
}

// GetFromServer mocks base method
func (m *MockQordobaClient) GetFromServer(getURL string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFromServer", getURL)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFromServer indicates an expected call of GetFromServer
func (mr *MockQordobaClientMockRecorder) GetFromServer(getURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFromServer", reflect.TypeOf((*MockQordobaClient)(nil).GetFromServer), getURL)
}

// PostToServer mocks base method
func (m *MockQordobaClient) PostToServer(postURL string, requestBody interface{}) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostToServer", postURL, requestBody)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostToServer indicates an expected call of PostToServer
func (mr *MockQordobaClientMockRecorder) PostToServer(postURL, requestBody interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostToServer", reflect.TypeOf((*MockQordobaClient)(nil).PostToServer), postURL, requestBody)
}

// PutToServer mocks base method
func (m *MockQordobaClient) PutToServer(postURL string, requestBody interface{}) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutToServer", postURL, requestBody)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutToServer indicates an expected call of PutToServer
func (mr *MockQordobaClientMockRecorder) PutToServer(postURL, requestBody interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutToServer", reflect.TypeOf((*MockQordobaClient)(nil).PutToServer), postURL, requestBody)
}

// DeleteFromServer mocks base method
func (m *MockQordobaClient) DeleteFromServer(deleteURL string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFromServer", deleteURL)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFromServer indicates an expected call of DeleteFromServer
func (mr *MockQordobaClientMockRecorder) DeleteFromServer(deleteURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFromServer", reflect.TypeOf((*MockQordobaClient)(nil).DeleteFromServer), deleteURL)
}

// MockLocal is a mock of Local interface
type MockLocal struct {
	ctrl     *gomock.Controller
	recorder *MockLocalMockRecorder
}

// MockLocalMockRecorder is the mock recorder for MockLocal
type MockLocalMockRecorder struct {
	mock *MockLocal
}

// NewMockLocal creates a new mock instance
func NewMockLocal(ctrl *gomock.Controller) *MockLocal {
	mock := &MockLocal{ctrl: ctrl}
	mock.recorder = &MockLocalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLocal) EXPECT() *MockLocalMockRecorder {
	return m.recorder
}

// Read mocks base method
func (m *MockLocal) Read(path string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", path)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockLocalMockRecorder) Read(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockLocal)(nil).Read), path)
}

// Write mocks base method
func (m *MockLocal) Write(fileName string, fileBytesResponse []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Write", fileName, fileBytesResponse)
}

// Write indicates an expected call of Write
func (mr *MockLocalMockRecorder) Write(fileName, fileBytesResponse interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockLocal)(nil).Write), fileName, fileBytesResponse)
}

// BuildFileName mocks base method
func (m *MockLocal) BuildFileName(file *types.File, filePathPattern, suffix string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildFileName", file, filePathPattern, suffix)
	ret0, _ := ret[0].(string)
	return ret0
}

// BuildFileName indicates an expected call of BuildFileName
func (mr *MockLocalMockRecorder) BuildFileName(file, suffix interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildFileName", reflect.TypeOf((*MockLocal)(nil).BuildFileName), file, suffix)
}

// FileExists mocks base method
func (m *MockLocal) FileExists(path string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileExists", path)
	ret0, _ := ret[0].(bool)
	return ret0
}

// FileExists indicates an expected call of FileExists
func (mr *MockLocalMockRecorder) FileExists(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileExists", reflect.TypeOf((*MockLocal)(nil).FileExists), path)
}

// QordobaHome mocks base method
func (m *MockLocal) QordobaHome() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QordobaHome")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QordobaHome indicates an expected call of QordobaHome
func (mr *MockLocalMockRecorder) QordobaHome() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QordobaHome", reflect.TypeOf((*MockLocal)(nil).QordobaHome))
}

// PutInHome mocks base method
func (m *MockLocal) PutInHome(fileName string, body []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PutInHome", fileName, body)
}

// PutInHome indicates an expected call of PutInHome
func (mr *MockLocalMockRecorder) PutInHome(fileName, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutInHome", reflect.TypeOf((*MockLocal)(nil).PutInHome), fileName, body)
}

// LoadCached mocks base method
func (m *MockLocal) LoadCached(cachedFileName string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadCached", cachedFileName)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadCached indicates an expected call of LoadCached
func (mr *MockLocalMockRecorder) LoadCached(cachedFileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadCached", reflect.TypeOf((*MockLocal)(nil).LoadCached), cachedFileName)
}

// FilesInFolder mocks base method
func (m *MockLocal) FilesInFolder(folderPath string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilesInFolder", folderPath)
	ret0, _ := ret[0].([]string)
	return ret0
}

// FilesInFolder indicates an expected call of FilesInFolder
func (mr *MockLocalMockRecorder) FilesInFolder(folderPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilesInFolder", reflect.TypeOf((*MockLocal)(nil).FilesInFolder), folderPath)
}

// RenderTable2Stdin mocks base method
func (m *MockLocal) RenderTable2Stdin(header []string, data [][]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderTable2Stdin", header, data)
}

// RenderTable2Stdin indicates an expected call of RenderTable2Stdin
func (mr *MockLocalMockRecorder) RenderTable2Stdin(header, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderTable2Stdin", reflect.TypeOf((*MockLocal)(nil).RenderTable2Stdin), header, data)
}

// MockConfigurationService is a mock of ConfigurationService interface
type MockConfigurationService struct {
	ctrl     *gomock.Controller
	recorder *MockConfigurationServiceMockRecorder
}

// MockConfigurationServiceMockRecorder is the mock recorder for MockConfigurationService
type MockConfigurationServiceMockRecorder struct {
	mock *MockConfigurationService
}

// NewMockConfigurationService creates a new mock instance
func NewMockConfigurationService(ctrl *gomock.Controller) *MockConfigurationService {
	mock := &MockConfigurationService{ctrl: ctrl}
	mock.recorder = &MockConfigurationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigurationService) EXPECT() *MockConfigurationServiceMockRecorder {
	return m.recorder
}

// ReadConfigInPath mocks base method
func (m *MockConfigurationService) ReadConfigInPath(path string) (*types.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadConfigInPath", path)
	ret0, _ := ret[0].(*types.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadConfigInPath indicates an expected call of ReadConfigInPath
func (mr *MockConfigurationServiceMockRecorder) ReadConfigInPath(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadConfigInPath", reflect.TypeOf((*MockConfigurationService)(nil).ReadConfigInPath), path)
}

// LoadConfig mocks base method
func (m *MockConfigurationService) LoadConfig() (*types.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadConfig")
	ret0, _ := ret[0].(*types.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadConfig indicates an expected call of LoadConfig
func (mr *MockConfigurationServiceMockRecorder) LoadConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadConfig", reflect.TypeOf((*MockConfigurationService)(nil).LoadConfig))
}

// SaveMainConfig mocks base method
func (m *MockConfigurationService) SaveMainConfig(config *types.Config) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveMainConfig", config)
}

// SaveMainConfig indicates an expected call of SaveMainConfig
func (mr *MockConfigurationServiceMockRecorder) SaveMainConfig(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMainConfig", reflect.TypeOf((*MockConfigurationService)(nil).SaveMainConfig), config)
}

// MockWorkspaceService is a mock of WorkspaceService interface
type MockWorkspaceService struct {
	ctrl     *gomock.Controller
	recorder *MockWorkspaceServiceMockRecorder
}

// MockWorkspaceServiceMockRecorder is the mock recorder for MockWorkspaceService
type MockWorkspaceServiceMockRecorder struct {
	mock *MockWorkspaceService
}

// NewMockWorkspaceService creates a new mock instance
func NewMockWorkspaceService(ctrl *gomock.Controller) *MockWorkspaceService {
	mock := &MockWorkspaceService{ctrl: ctrl}
	mock.recorder = &MockWorkspaceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWorkspaceService) EXPECT() *MockWorkspaceServiceMockRecorder {
	return m.recorder
}

// LoadWorkspace mocks base method
func (m *MockWorkspaceService) LoadWorkspace() (*types.WorkspaceData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadWorkspace")
	ret0, _ := ret[0].(*types.WorkspaceData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadWorkspace indicates an expected call of LoadWorkspace
func (mr *MockWorkspaceServiceMockRecorder) LoadWorkspace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadWorkspace", reflect.TypeOf((*MockWorkspaceService)(nil).LoadWorkspace))
}

// MockFileService is a mock of FileService interface
type MockFileService struct {
	ctrl     *gomock.Controller
	recorder *MockFileServiceMockRecorder
}

// MockFileServiceMockRecorder is the mock recorder for MockFileService
type MockFileServiceMockRecorder struct {
	mock *MockFileService
}

// NewMockFileService creates a new mock instance
func NewMockFileService(ctrl *gomock.Controller) *MockFileService {
	mock := &MockFileService{ctrl: ctrl}
	mock.recorder = &MockFileServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileService) EXPECT() *MockFileServiceMockRecorder {
	return m.recorder
}

// WorkspaceFiles mocks base method
func (m *MockFileService) WorkspaceFiles(personaID int, withProgressStatus bool) (*types.FileSearchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WorkspaceFiles", personaID, withProgressStatus)
	ret0, _ := ret[0].(*types.FileSearchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WorkspaceFiles indicates an expected call of WorkspaceFiles
func (mr *MockFileServiceMockRecorder) WorkspaceFiles(personaID, withProgressStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WorkspaceFiles", reflect.TypeOf((*MockFileService)(nil).WorkspaceFiles), personaID, withProgressStatus)
}

// WorkspaceFilesWithLimit mocks base method
func (m *MockFileService) WorkspaceFilesWithLimit(personaID int, withProgressStatus bool, limit int) (*types.FileSearchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WorkspaceFilesWithLimit", personaID, withProgressStatus, limit)
	ret0, _ := ret[0].(*types.FileSearchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WorkspaceFilesWithLimit indicates an expected call of WorkspaceFilesWithLimit
func (mr *MockFileServiceMockRecorder) WorkspaceFilesWithLimit(personaID, withProgressStatus, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WorkspaceFilesWithLimit", reflect.TypeOf((*MockFileService)(nil).WorkspaceFilesWithLimit), personaID, withProgressStatus, limit)
}

// FindFile mocks base method
func (m *MockFileService) FindFile(fileName, version string, withProgressStatus bool) (*types.File, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindFile", fileName, version, withProgressStatus)
	ret0, _ := ret[0].(*types.File)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// FindFile indicates an expected call of FindFile
func (mr *MockFileServiceMockRecorder) FindFile(fileName, version, withProgressStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFile", reflect.TypeOf((*MockFileService)(nil).FindFile), fileName, version, withProgressStatus)
}

// DownloadFile mocks base method
func (m *MockFileService) DownloadFile(personaID int, fileName string, file *types.File) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DownloadFile", personaID, fileName, file)
}

// DownloadFile indicates an expected call of DownloadFile
func (mr *MockFileServiceMockRecorder) DownloadFile(personaID, fileName, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFile", reflect.TypeOf((*MockFileService)(nil).DownloadFile), personaID, fileName, file)
}

// DownloadSourceFile mocks base method
func (m *MockFileService) DownloadSourceFile(fileName string, file *types.File, withUpdates bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DownloadSourceFile", fileName, file, withUpdates)
}

// DownloadSourceFile indicates an expected call of DownloadSourceFile
func (mr *MockFileServiceMockRecorder) DownloadSourceFile(fileName, file, withUpdates interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadSourceFile", reflect.TypeOf((*MockFileService)(nil).DownloadSourceFile), fileName, file, withUpdates)
}

// PushFolder mocks base method
func (m *MockFileService) PushFolder(folder, version string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PushFolder", folder, version)
}

// PushFolder indicates an expected call of PushFolder
func (mr *MockFileServiceMockRecorder) PushFolder(folder, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushFolder", reflect.TypeOf((*MockFileService)(nil).PushFolder), folder, version)
}

// PushFiles mocks base method
func (m *MockFileService) PushFiles(fileList []string, version string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PushFiles", fileList, version)
}

// PushFiles indicates an expected call of PushFiles
func (mr *MockFileServiceMockRecorder) PushFiles(fileList, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushFiles", reflect.TypeOf((*MockFileService)(nil).PushFiles), fileList, version)
}

// DeleteFile mocks base method
func (m *MockFileService) DeleteFile(fileName, version string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteFile", fileName, version)
}

// DeleteFile indicates an expected call of DeleteFile
func (mr *MockFileServiceMockRecorder) DeleteFile(fileName, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockFileService)(nil).DeleteFile), fileName, version)
}

// FileScore mocks base method
func (m *MockFileService) FileScore(filename, version string) *types.ScoreResponseBody {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileScore", filename, version)
	ret0, _ := ret[0].(*types.ScoreResponseBody)
	return ret0
}

// FileScore indicates an expected call of FileScore
func (mr *MockFileServiceMockRecorder) FileScore(filename, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileScore", reflect.TypeOf((*MockFileService)(nil).FileScore), filename, version)
}

// MockSegmentService is a mock of SegmentService interface
type MockSegmentService struct {
	ctrl     *gomock.Controller
	recorder *MockSegmentServiceMockRecorder
}

// MockSegmentServiceMockRecorder is the mock recorder for MockSegmentService
type MockSegmentServiceMockRecorder struct {
	mock *MockSegmentService
}

// NewMockSegmentService creates a new mock instance
func NewMockSegmentService(ctrl *gomock.Controller) *MockSegmentService {
	mock := &MockSegmentService{ctrl: ctrl}
	mock.recorder = &MockSegmentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSegmentService) EXPECT() *MockSegmentServiceMockRecorder {
	return m.recorder
}

// FindSegment mocks base method
func (m *MockSegmentService) FindSegment(fileName, fileVersion, key string) (*types.Segment, *types.File) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSegment", fileName, fileVersion, key)
	ret0, _ := ret[0].(*types.Segment)
	ret1, _ := ret[1].(*types.File)
	return ret0, ret1
}

// FindSegment indicates an expected call of FindSegment
func (mr *MockSegmentServiceMockRecorder) FindSegment(fileName, fileVersion, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSegment", reflect.TypeOf((*MockSegmentService)(nil).FindSegment), fileName, fileVersion,
		key)
}

// AddKey mocks base method
func (m *MockSegmentService) AddKey(fileName, version string, keyAddRequest *types.KeyAddRequest) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddKey", fileName, version, keyAddRequest)
}

// AddKey indicates an expected call of AddKey
func (mr *MockSegmentServiceMockRecorder) AddKey(fileName, version, keyAddRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddKey", reflect.TypeOf((*MockSegmentService)(nil).AddKey), fileName, version, keyAddRequest)

}

// UpdateKey mocks base method
func (m *MockSegmentService) UpdateKey(fileName, version string, keyAddRequest *types.KeyAddRequest) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateKey", fileName, version, keyAddRequest)
}

// UpdateKey indicates an expected call of UpdateKey
func (mr *MockSegmentServiceMockRecorder) UpdateKey(fileName, version, keyAddRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateKey", reflect.TypeOf((*MockSegmentService)(nil).UpdateKey), fileName, version, keyAddRequest)
}

// DeleteKey mocks base method
func (m *MockSegmentService) DeleteKey(fileName, version, segmentKey string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteKey", fileName, version, segmentKey)
}

// DeleteKey indicates an expected call of DeleteKey
func (mr *MockSegmentServiceMockRecorder) DeleteKey(fileName, version, segmentKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteKey", reflect.TypeOf((*MockSegmentService)(nil).DeleteKey), fileName, version, segmentKey)
}
