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
	"time"
)

const (
	pushFileTemplate = "%s/v3/files/organizations/%d/workspaces/%d/upsert"
	concurrencyLevel = 1
)

func (f *FileService) PushFolder(folder, version string) {
	fileList := f.Local.FilesInFolder(folder)
	f.PushFiles(fileList, version)
}

func (f *FileService) PushFiles(fileList []string, version string) {
	jobs := make(chan *pushFileTask, 1000)
	results := make(chan struct{}, 1000)
	filteredFileList := f.filterFiles(fileList)

	for i := 0; i < concurrencyLevel; i++ {
		go f.startPushWorker(jobs, results, version)
	}

	// let all error logs go before final messages
	time.Sleep(time.Second)
	totalFilesPushed := 0
	for _, file := range filteredFileList {
		totalFilesPushed += f.pushFile(file, jobs)
	}
	close(jobs)
	for i := 0; i < totalFilesPushed; i++ {
		<-results
	}
}

func (f *FileService) filterFiles(files []string) []string {
	filteredFiles := make([]string, 0, 0)
	blacklistRegexp := make([]*regexp.Regexp, 0, len(f.Config.Blacklist.Sources))
	for _, blackList := range f.Config.Blacklist.Sources {
		compile, err := regexp.Compile(blackList)
		if err != nil {
			log.Errorf("invalid blacklist regexp '%s': %v\n", blackList, err)
			os.Exit(1)
		}
		blacklistRegexp = append(blacklistRegexp, compile)
	}
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

func (f *FileService) startPushWorker(jobs chan *pushFileTask, results chan struct{}, version string) {
	base := f.Config.GetAPIBase()
	pushFileURL := fmt.Sprintf(pushFileTemplate, base, f.Config.Qordoba.OrganizationID, f.Config.Qordoba.WorkspaceID)
	for j := range jobs {
		f.sendFileToServer(j.fileInfo, j.FilePath, pushFileURL, version)
		results <- struct{}{}
	}
}

func (f *FileService) pushFile(filePath string, jobs chan *pushFileTask) int {
	fileInfo, e := os.Stat(filePath)
	if e != nil {
		log.Errorf("error occurred in file read: %v", e)
		return 0
	}
	if fileInfo.IsDir() {
		return f.handleDirectory2Push(filePath, jobs)
	}
	jobs <- &pushFileTask{
		FilePath: filePath,
		fileInfo: fileInfo,
	}
	return 1
}

func (f *FileService) handleDirectory2Push(filePath string, jobs chan *pushFileTask) int {
	file2Push := 0
	err := filepath.Walk(filePath, func(path string, childFileInfo os.FileInfo, err error) error {
		if !childFileInfo.IsDir() {
			jobs <- &pushFileTask{
				FilePath: filePath,
				fileInfo: childFileInfo,
			}
			file2Push++
		}
		return nil
	})
	if err != nil {
		log.Errorf("error occurred: %v", err)
	}
	return file2Push
}

type pushFileTask struct {
	FilePath string
	fileInfo os.FileInfo
}

func (f *FileService) sendFileToServer(fileInfo os.FileInfo, filePath, pushFileURL, version string) {
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
		log.Errorf("error occurred on building file post request: %v", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		if resp.StatusCode == http.StatusUnauthorized {
			log.Errorf("User is not authorised for this request. Check `access_token` in configuration.")
		} else {
			log.Errorf("File %s push status: %v. Response : %v", filePath, resp.Status, string(body))
		}
	} else {
		if version == "" {
			log.Infof("File %s was pushed to server.", filePath)
		} else {
			log.Infof("File %s (version '%v') was pushed to server.", filePath, version)
		}
	}
}

func (f *FileService) buildPushRequest(fileInfo os.FileInfo, filePath, version string) (*types.PushRequest, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("can't handle file %s: %v", filePath, err)
		return nil, err
	}
	return &types.PushRequest{
		FileName: fileInfo.Name(),
		Version:  version,
		Content:  string(fileContent),
	}, nil
}
