package file

import (
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
	"time"
)

const (
	pushFileTemplate = "%s/v3/files/organizations/%d/workspaces/%d/upsert"
	concurrencyLevel = 1
)

// PushFolder function push folder to server
func (f *Service) PushFolder(folder, version string) {
	fileList := f.Local.FilesInFolder(folder)
	f.PushFiles(fileList, version)
}

// PushFiles function push array of files to server with specified version
func (f *Service) PushFiles(fileList []string, version string) {
	jobs := make(chan *pushFileTask, 1000)
	results := make(chan struct{}, 1000)
	filteredFileList := f.filterFiles(fileList)

	for i := 0; i < concurrencyLevel; i++ {
		go f.startPushWorker(jobs, results, version)
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
		if f.Config.Push.LanguageCode != "" && !strings.Contains(file, f.Config.Push.LanguageCode) {
			log.Infof("file %s is not pushed due to doesn't contain %s langage code in path", file, f.Config.Push.LanguageCode)
			continue
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

func (f *Service) startPushWorker(jobs chan *pushFileTask, results chan struct{}, version string) {
	base := f.Config.GetAPIBase()
	pushFileURL := fmt.Sprintf(pushFileTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.WorkspaceID)
	for j := range jobs {
		f.sendFileToServer(j.fileInfo, j.FilePath, pushFileURL, version, results)
	}
}

func (f *Service) pushFile(filePath string, jobs chan *pushFileTask) int {
	fileInfo, e := os.Stat(filePath)
	if e != nil {
		log.Errorf("error occurred in file read: %v", e)
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

func (f *Service) sendFileToServer(fileInfo os.FileInfo, filePath, pushFileURL, version string, results chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Recovered in sendFileToServer: %v\n%s\n", err, debug.Stack())
		}
		results <- struct{}{}
	}()
	if fileInfo.IsDir() {
		// this is possible in case of folder presence in folder. Currently we don't support recursion, so just ignore
		return
	}
	pushRequest, err := f.buildPushRequest(fileInfo, filePath, version)
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

func (f *Service) buildPushRequest(fileInfo os.FileInfo, filePath, version string) (*types.PushRequest, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("can't handle file %s: %v", filePath, err)
		return nil, err
	}
	dir, _ := os.Getwd()
	relativeFilePath, _ := filepath.Rel(dir, filePath)
	relativeFilePath = strings.ReplaceAll(relativeFilePath, "\\", "/")
	return &types.PushRequest{
		FileName: fileInfo.Name(),
		Version:  version,
		Content:  string(fileContent),
		Filepath: relativeFilePath,
	}, nil
}
