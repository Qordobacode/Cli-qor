package file

import (
	"errors"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"
)

const (
	pushFileTemplate = "%s/v3/files/organizations/%d/workspaces/%d/upsert"
	concurrencyLevel = 1
)

var (
	TotalSkipped uint64
	MimeTypes    string
)

// PushFolder function push folder to server
func (f *Service) PushFolder(folder, version string, isRecursive bool) {
	fileList := f.Local.FilesInFolder(folder, isRecursive)
	f.PushFiles(fileList, version, isRecursive)
}

// PushFiles function push array of files to server with specified version
func (f *Service) PushFiles(fileList []string, version string, isFilepath bool) {
	jobs := make(chan *pushFileTask, 1000)
	results := make(chan struct{}, 1000)
	filteredFileList := f.filterFiles(fileList)

	workspace, err := f.WorkspaceService.LoadWorkspace()
	if err != nil {
		os.Exit(1)
	}
	contentTypeCodes := map[string]struct{}{}
	contentTypeArray := make([]string, 0)
	for _, code := range workspace.Workspace.ContentTypeCodes {
		for _, ext := range code.Extensions {
			contentTypeCodes[ext] = struct{}{}
			contentTypeArray = append(contentTypeArray, `"`+ext+`"`)
		}
	}
	MimeTypes = strings.Join(contentTypeArray, ", ")
	for i := 0; i < concurrencyLevel; i++ {
		go f.startPushWorker(jobs, results, version, workspace, contentTypeCodes, isFilepath)
	}

	// let all error logs go before final messages
	time.Sleep(time.Second)
	totalFilesPushed := 0
	for i := range filteredFileList {
		totalFilesPushed += f.pushFile(filteredFileList[i], jobs)
	}
	close(jobs)
	for i := 0; i < totalFilesPushed; i++ {
		<-results
	}
}

func (f *Service) filterFiles(files []string) []string {
	filteredFiles := make([]string, 0, 0)
	blacklistRegexp := f.buildBlacklistRegexps()

fileSearch:
	for _, file := range files {
		for _, blackReg := range blacklistRegexp {
			if blackReg.FindString(file) != "" {
				log.Infof("file %s is not pushed due to black list", file)
				continue fileSearch
			}
		}
		filteredFiles = append(filteredFiles, file)
	}
	return filteredFiles
}

func (f *Service) buildBlacklistRegexps() []*regexp.Regexp {
	blacklistRegexp := make([]*regexp.Regexp, 0, len(f.Config.Blacklist.Sources))
	for _, blackList := range f.Config.Blacklist.Sources {
		compile, err := regexp.Compile(blackList)
		if err != nil {
			log.Errorf("invalid blacklist regexp '%s': %v\n", blackList, err)
			os.Exit(1)
		}
		blacklistRegexp = append(blacklistRegexp, compile)
	}
	return blacklistRegexp
}

func (f *Service) startPushWorker(jobs chan *pushFileTask, results chan struct{}, version string,
	workspace *types.WorkspaceData, contentTypeCodes map[string]struct{}, isFilepath bool) {
	base := f.Config.GetAPIBase()
	pushFileURL := fmt.Sprintf(pushFileTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.WorkspaceID)
	for j := range jobs {
		f.sendFileToServer(j.fileInfo, j.FilePath, pushFileURL, version, results, workspace, contentTypeCodes, isFilepath)
	}
}

func (f *Service) pushFile(filePath string, jobs chan *pushFileTask) int {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Errorf("file %s doesn't exist", filePath)
		return 0
	}
	if err != nil {
		log.Errorf("error occurred on file read: %v", err)
		return 0
	}
	if fileInfo.IsDir() {
		return 0
	}
	jobs <- &pushFileTask{
		FilePath: filePath,
		fileInfo: fileInfo,
	}
	return 1
}

type pushFileTask struct {
	FilePath string
	fileInfo os.FileInfo
}

func (f *Service) sendFileToServer(fileInfo os.FileInfo, filePath, pushFileURL, version string, results chan struct{},
	workspace *types.WorkspaceData, contentTypeCodes map[string]struct{}, isFilepath bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Recovered in sendFileToServer: %v\n%s\n", err, debug.Stack())
		}
		results <- struct{}{}
	}()
	if fileInfo.IsDir() {
		return
	}
	pushRequest, err := f.buildPushRequest(fileInfo, filePath, version, workspace, contentTypeCodes, isFilepath)
	if err != nil {
		return
	}
	resp, err := f.QordobaClient.PostToServer(pushFileURL, pushRequest)
	if err != nil {
		log.Errorf("error occurred on post to server: %v", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		if resp.StatusCode == http.StatusUnauthorized {
			log.Errorf("User is not authorised for this request. Check `access_token` in configuration.")
			os.Exit(1)
		}
		if resp.StatusCode == http.StatusRequestEntityTooLarge {
			log.Errorf("File %v (%v bytes) is too large for server. %v", fileInfo.Name(), fileInfo.Size(), string(body))
		} else {
			log.Errorf("File %s push status: %v. Response: %v", filePath, resp.Status, string(body))
		}
	} else {
		if version == "" {
			log.Infof("File %s was pushed to server.", filePath)
		} else {
			log.Infof("File %s (version '%v') was pushed to server.", filePath, version)
		}
	}
}

func (f *Service) buildPushRequest(fileInfo os.FileInfo, filePath, version string, workspace *types.WorkspaceData,
	contentTypeCodes map[string]struct{}, isFilepath bool) (*types.PushRequest, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("can't handle file %s: %v", filePath, err)
		return nil, err
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Debugf("error occurred on getting current dir: %v", err)
	}
	if len(f.Config.Push.Sources.Folders) > 0 {
		dir = f.Config.Push.Sources.Folders[0]
	}
	relativeFilePath, err := filepath.Rel(dir, filePath)
	if err != nil {
		log.Debugf("error occurred on building relativePath: %v", err)
		relativeFilePath = filePath
	}
	if relativeFilePath == "" {
		log.Debugf("relativeFilePath is empty. Use filePath '%s' instead", filePath)
		relativeFilePath = filePath
	}
	if isFilepath && !filterFileByWorkspace(relativeFilePath, filePath, workspace) {
		return nil, errors.New("file not pass source name")
	}
	if !filterFileByMimeType(filePath, fileInfo.Name(), contentTypeCodes) {
		return nil, errors.New("file not pass mime type check")
	}
	if !isFilepath {
		relativeFilePath = ""
	}
	relativeFilePath = filepath.Dir(relativeFilePath)
	relativeFilePath = strings.ReplaceAll(relativeFilePath, "\\", "/")
	return &types.PushRequest{
		FileName: fileInfo.Name(),
		Version:  version,
		Content:  string(fileContent),
		Filepath: relativeFilePath,
	}, nil
}

func filterFileByMimeType(filePath, fileName string, contentTypeCodes map[string]struct{}) bool {
	if len(contentTypeCodes) == 0 {
		// pass checks if workspace's content type is empty
		return true
	}
	dotParts := strings.Split(fileName, ".")
	if len(dotParts) < 2 {
		log.Infof("[SKIPPED] File '%s' doesn't contain workspace mime type", filePath)
		atomic.AddUint64(&TotalSkipped, 1)
		return false
	}
	mimeType := dotParts[len(dotParts)-1]
	if _, ok := contentTypeCodes[mimeType]; !ok {
		log.Infof("[SKIPPED] File '%s' doesn't contain workspace mime type", filePath)
		atomic.AddUint64(&TotalSkipped, 1)
		return false
	}
	return true
}

func filterFileByWorkspace(relativeFilePath, filePath string, workspace *types.WorkspaceData) bool {
	log.Debugf("filepath = '%s', relativeFilePath = '%s'", filePath, relativeFilePath)
	relativeFilePath = strings.ToLower(relativeFilePath)
	filePaths := strings.Split(relativeFilePath, ".")
	noMimeRelativeFilePath := strings.TrimSuffix(relativeFilePath, "."+filePaths[len(filePaths)-1])
	code := workspace.Workspace.SourcePersona.Code
	codeSplits := strings.Split(code, "-")
	nameSplits := strings.Split(workspace.Workspace.SourcePersona.Name, "-")
	for _, codeVal := range codeSplits {
		codeVal = strings.TrimSpace(strings.ToLower(codeVal))
		if strings.Contains(noMimeRelativeFilePath, string(os.PathSeparator)+codeVal+"-") ||
			strings.Contains(noMimeRelativeFilePath, codeVal+string(os.PathSeparator)) ||
			strings.HasPrefix(noMimeRelativeFilePath, codeVal+"-") ||
			strings.HasSuffix(noMimeRelativeFilePath, codeVal) {
			return true
		}
		log.Debugf("relativeFilePath = '%s' doesn't contain code '%s'", relativeFilePath, codeVal)
	}
	for _, name := range nameSplits {
		name = strings.TrimSpace(strings.ToLower(name))
		if strings.Contains(noMimeRelativeFilePath, string(os.PathSeparator)+name+"-") ||
			strings.Contains(noMimeRelativeFilePath, name+string(os.PathSeparator)) ||
			strings.HasPrefix(noMimeRelativeFilePath, name+"-") ||
			strings.HasSuffix(noMimeRelativeFilePath, name) {
			return true
		}
		log.Debugf("relativeFilePath = '%s' doesn't contain name '%s'", relativeFilePath, name)
	}

	folderErrMsg := ""
	fileErrMsg := ""
	if len(codeSplits) > 1 {
		folderErrMsg = fmt.Sprintf(`/%s/, /%s-xxx/, /xxx-%s/, /%s/, /%s-xxx/, /xxx-%s/`,
			codeSplits[0], codeSplits[0], codeSplits[0], code, code, code)
		fileErrMsg = fmt.Sprintf(`%s., -%s., %s-, %s., %s-, -%s`,
			codeSplits[0], codeSplits[0], codeSplits[0], code, code, code)
	}
	log.Infof("[SKIPPED] File path '%s' doesn't contain Source code. For example: %s or in file name %s.", filePath, folderErrMsg, fileErrMsg)
	return false
}
