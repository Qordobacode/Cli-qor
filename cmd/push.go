package cmd

import (
	"bytes"
	"fmt"
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

const (
	pushFileTemplate    = "%s/v3/files/organizations/%d/workspaces/%d/upsert"
	ApplicationJsonType = "application/json"
)

var (
	fileVersion         string
	isUpdate            bool
	allowedMimeTypes, _ = regexp.Compile("\\.(csv|xml|json|txt|yaml|yml)$")

	// HTTPClient - custom one with a delay set
	HTTPClient = http.Client{
		Timeout: time.Minute * 10,
	}
	// pushCmd represents the push command
	pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push files or folders",
		Run:   pushCommand,
	}
)

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

func init() {
	pushCmd.Flags().StringVarP(&fileVersion, "version", "v", "", "--version")
	pushCmd.Flags().BoolVarP(&isUpdate, "update", "u", false, "--update")
	rootCmd.AddCommand(pushCmd)
}

func getFolderFileNames() []string {
	result := make([]string, 0, 0)
	curFolderFiles, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Printf("error occurred on retrieving list of all files in current folder: %v\n", err)
		return result
	}
	for _, f := range curFolderFiles {
		file := f.Name()
		if allowedMimeTypes.FindString(file) != "" {
			fmt.Printf("push file: %v\n", file)
		}
		result = append(result, file)
	}
	return result
}

func pushFile(qordoba *general.QordobaConfig, filePath string) {
	pathContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Infof("can't handle file %s: %v\n", filePath, err)
		return
	}
	base := qordoba.GetAPIBase()
	pushFileURL := fmt.Sprintf(pushFileTemplate, base, qordoba.Qordoba.OrganizationID, qordoba.Qordoba.ProjectID)
	reader := bytes.NewReader(pathContent)
	request, err := http.NewRequest("POST", pushFileURL, reader)
	if err != nil {
		log.Infof("error occurred on building file post request: %v\n", err)
		return
	}
	request.Header.Add("x-auth-token", qordoba.Qordoba.AccessToken)
	request.Header.Add("Content-Type", ApplicationJsonType)
	_, err = HTTPClient.Do(request)
	if err != nil {
		log.Infof("error occurred on sending POST request to server")
	}
}
