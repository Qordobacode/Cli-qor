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
	"os"
	"path/filepath"
	"regexp"
)

const (
	pushFileTemplate = "%s/v3/files/organizations/%d/workspaces/%d/upsert"
)

var (
	fileVersion         string
	isUpdate            bool
	allowedMimeTypes, _ = regexp.Compile("\\.(csv|xml|json|txt|yaml|yml)$")

	// pushCmd represents the push command
	pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push files or folders",
		Run:   pushCommand,
	}
)

func init() {
	pushCmd.Flags().StringVarP(&fileVersion, "version", "v", "", "--version")
	pushCmd.Flags().BoolVarP(&isUpdate, "update", "u", false, "--update")
	rootCmd.AddCommand(pushCmd)
}

func pushCommand(cmd *cobra.Command, args []string) {
	log.Debugf("push was called\n")
	log.Debugf("version = %v\n", fileVersion)
	log.Debugf("args: %v\n", args)
	qordobaConfig, err := general.LoadConfig()
	if err != nil {
		return
	}
	fileList := args
	if len(args) == 0 {
		fileList = getFolderFileNames()
	}
	for _, file := range fileList {
		pushFile(qordobaConfig, file)
	}
}

func getFolderFileNames() []string {
	result := make([]string, 0, 0)
	curFolderFiles, err := ioutil.ReadDir("./")
	if err != nil {
		log.Errorf("error occurred on retrieving list of all files in current folder: %v\n", err)
		return result
	}
	for _, f := range curFolderFiles {
		file := f.Name()
		if allowedMimeTypes.FindString(file) != "" {
			//fmt.Printf("push file: %v\n", file)
		}
		result = append(result, file)
	}
	return result
}

func pushFile(qordoba *general.QordobaConfig, filePath string) {
	fileInfo, e := os.Stat(filePath)
	if e != nil {
		log.Errorf("error occurred in file read: %v\n", e)
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

func sendFileToServer(fileInfo os.FileInfo, qordoba *general.QordobaConfig, filePath, pushFileURL string) {
	if fileInfo.IsDir() {
		// this is possible in case of folder presence in folder. Currently we don't support recursion, so just ignore
		return
	}
	reader, err := buildPushRequestBody(fileInfo, filePath)
	if err != nil {
		return
	}
	general.PostToServer(qordoba, filePath, pushFileURL, reader)
}

func buildPushRequestBody(fileInfo os.FileInfo, filePath string) (io.Reader, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("can't handle file %s: %v\n", filePath, err)
		return nil, err
	}
	requestBody := general.PushRequest{
		FileName: fileInfo.Name(),
		Version:  fileVersion,
		Content:  string(fileContent),
	}

	marshaledBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Errorf("error occurred on marshalling object: %v\n", err)
		return nil, err
	}
	reader := bytes.NewReader(marshaledBody)
	return reader, nil
}
