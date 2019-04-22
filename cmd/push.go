package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

const (
	pushFileTemplate = "%s/v3/files/organizations/%d/workspaces/%d/upsert"
)

var (
	tag                 string
	files               string
	folderPath          string
	allowedMimeTypes, _ = regexp.Compile("\\.(csv|xml|json|txt|yaml|yml)$")

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
		log.Infof("no '--files' or '--file-path' params in command. 'source' param from config is used\n  Files: %v\n  Folders: %v", pushSources.Files, pushSources.Folders)
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

func pushFiles(fileList []string, qordobaConfig *general.Config) {
	for _, file := range fileList {
		pushFile(qordobaConfig, file)
	}
}

func getFilesInFolder(folderPath string) []string {
	result := make([]string, 0, 0)
	curFolderFiles, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Errorf("error occurred on retrieving list of all files in current folder: %v", err)
		return result
	}
	for _, f := range curFolderFiles {
		file := folderPath + string(os.PathSeparator) + f.Name()
		if allowedMimeTypes.FindString(file) != "" {
			//fmt.Printf("push file: %v", file)
		}
		result = append(result, file)
	}
	return result
}

func pushFile(qordoba *general.Config, filePath string) {
	fileInfo, e := os.Stat(filePath)
	if e != nil {
		log.Errorf("error occurred in file read: %v", e)
		return
	}
	base := qordoba.GetAPIBase()
	pushFileURL := fmt.Sprintf(pushFileTemplate, base, qordoba.Qordoba.OrganizationID, qordoba.Qordoba.ProjectID)
	if fileInfo.IsDir() {
		err := filepath.Walk(filePath, func(path string, childFileInfo os.FileInfo, err error) error {
			if childFileInfo.IsDir() {
				return nil
			}
			sendFileToServer(childFileInfo, qordoba, path, pushFileURL)
			return nil
		})
		if err != nil {
			log.Errorf("error occurred: %v", err)
		}
	} else {
		sendFileToServer(fileInfo, qordoba, filePath, pushFileURL)
	}

}

func sendFileToServer(fileInfo os.FileInfo, qordoba *general.Config, filePath, pushFileURL string) {
	if fileInfo.IsDir() {
		// this is possible in case of folder presence in folder. Currently we don't support recursion, so just ignore
		return
	}
	reader, err := buildPushRequestBody(fileInfo, filePath)
	if err != nil {
		return
	}
	resp, err := general.PostToServer(qordoba, pushFileURL, reader)
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
			log.Infof("File %s was pushed to server.", filePath)
		} else {
			log.Infof("File %s (version '%v') was pushed to server.", filePath, tag)
		}
	}
}

func buildPushRequestBody(fileInfo os.FileInfo, filePath string) (io.Reader, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("can't handle file %s: %v", filePath, err)
		return nil, err
	}
	requestBody := general.PushRequest{
		FileName: fileInfo.Name(),
		Version:  tag,
		Content:  string(fileContent),
	}

	marshaledBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Errorf("error occurred on marshalling object: %v", err)
		return nil, err
	}
	reader := bytes.NewReader(marshaledBody)
	return reader, nil
}
