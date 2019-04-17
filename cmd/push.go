package cmd

import (
	"fmt"
	"github.com/qordobacode/cli-v2/files"
	"github.com/qordobacode/cli-v2/log"
	"github.com/qordobacode/cli-v2/models"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var (
	fileVersion         string
	isUpdate            bool
	allowedMimeTypes, _ = regexp.Compile("\\.(csv|xml|json|txt|yaml|yml)$")

	// HTTPClient - custom one with a delay set
	HTTPClient = http.Client{
		Timeout: time.Minute * 10,
	}
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push files or folders",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("push was called\n")
		log.Debugf("version = %v\n", fileVersion)
		log.Debugf("args: %v\n", args)
		qordobaConfig, err := files.LoadConfig()
		if err != nil {
			return
		}
		fileList := args
		if len(args) == 0 {
			fileList = getFolderFileNames()
		}
		pushFiles(qordobaConfig, fileList)
	},
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

func pushFiles(qordoba *models.QordobaConfig, files []string) {

}

func pushFile(filePath string) {
	bytes, err := ioutil.ReadFile(filePath)

}
