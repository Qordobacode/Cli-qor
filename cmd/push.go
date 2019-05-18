package cmd

import (
	"fmt"
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	pushFileTemplate = "%s/v3/files/organizations/%d/workspaces/%d/upsert"
)

var (
	tag                 string
	files               string
	folderPath          string
	allowedMimeTypes, _ = regexp.Compile(`\.(csv|xml|json|txt|yaml|yml)$`)

	// pushCmd represents the push command
	pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push files or folders",
		Run:   pushCommand,
	}
)

func init() {
	pushCmd.Flags().StringVarP(&tag, "version", "v", "", "--version")
	pushCmd.Flags().StringVarP(&files, "files", "f", "", "--update")
	pushCmd.Flags().StringVarP(&folderPath, "file-path", "p", "", "--update")
	rootCmd.AddCommand(pushCmd)
}

func pushCommand(cmd *cobra.Command, args []string) {
	qordobaConfig, err := general.LoadConfig()
	if err != nil {
		return
	}
	if folderPath == "" && files == "" {
		pushSources := qordobaConfig.Push.Sources
		log.Infof("no '--files' or '--file-path' params in command. 'source' param from config is used\n  File: %v\n  Folders: %v", pushSources.Files, pushSources.Folders)
		pushFiles(pushSources.Files, qordobaConfig)
		for _, folder := range pushSources.Folders {
			pushFolder(folder, qordobaConfig)
		}
	} else {
		if files != "" {
			fileList := filepath.SplitList(files)
			log.Debugf("Result list of files from line is: %v", string(os.PathListSeparator), fileList)
			pushFiles(fileList, qordobaConfig)
		}
		if folderPath != "" {
			pushFolder(folderPath, qordobaConfig)
		}
	}
}

func pushFolder(folder string, qordobaConfig *general.Config) {
	fileList := getFilesInFolder(folder)
	pushFiles(fileList, qordobaConfig)
}

func getFilesInFolder(folderPath string) []string {
	result := make([]string, 0)
	curFolderFiles, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Errorf("error occurred on retrieving list of all files in current folder: %v", err)
		return result
	}
	for _, f := range curFolderFiles {
		file := folderPath + string(os.PathSeparator) + f.Name()
		result = append(result, file)
	}
	return result
}

func pushFiles(fileList []string, config *general.Config) {
	jobs := make(chan *pushFileTask, 1000)
	results := make(chan struct{}, 1000)

	for i := 0; i < 3; i++ {
		go startPushWorker(config, jobs, results)
	}

	// let all error logs go before final messages
	time.Sleep(time.Second)
	totalFilesPushed := 0
	for _, file := range fileList {
		totalFilesPushed += pushFile(file, jobs)
	}
	close(jobs)
	for i:= 0; i < totalFilesPushed; i++ {
		<-results
	}
}

func startPushWorker(config *general.Config, jobs chan *pushFileTask, results chan struct{}) {
	base := config.GetAPIBase()
	pushFileURL := fmt.Sprintf(pushFileTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID)
	for j := range jobs {
		sendFileToServer(config, j.fileInfo, j.FilePath, pushFileURL)
		results <- struct{}{}
	}
}

func pushFile(filePath string, jobs chan *pushFileTask) int {
	fileInfo, e := os.Stat(filePath)
	if e != nil {
		log.Errorf("error occurred in file read: %v", e)
		return 0
	}
	if fileInfo.IsDir() {
		return handleDirectory2Push(filePath, jobs)
	}
	jobs <- &pushFileTask{
		FilePath: filePath,
		fileInfo: fileInfo,
	}
	return 1
}

func handleDirectory2Push(filePath string, jobs chan *pushFileTask) int {
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

func sendFileToServer(config *general.Config, fileInfo os.FileInfo, filePath, pushFileURL string) {
	if fileInfo.IsDir() {
		// this is possible in case of folder presence in folder. Currently we don't support recursion, so just ignore
		return
	}
	pushRequest, err := buildPushRequest(fileInfo, filePath)
	if err != nil {
		return
	}
	resp, err := general.PostToServer(config, pushFileURL, pushRequest)
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
		if tag == "" {
			log.Infof("File '%s' was pushed to server.", filePath)
		} else {
			log.Infof("File '%s' (version '%v') was pushed to server.", filePath, tag)
		}
	}
}

func buildPushRequest(fileInfo os.FileInfo, filePath string) (*general.PushRequest, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("can't handle file %s: %v", filePath, err)
		return nil, err
	}
	return &general.PushRequest{
		FileName: fileInfo.Name(),
		Version:  tag,
		Content:  string(fileContent),
	}, nil
}
