package cmd

import (
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"
	"github.com/spf13/cobra"
)

var (
	headers = []string{
		"#AUDIENSES", "#WORDS", "#SEGMENTS", "EDITING", "PROOFREADING", "COMPLETED",
	}
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status per project or file (Support file versions)",
	Run:   runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	config, err := general.LoadConfig()
	if err != nil {
		log.Errorf("error occurred on configuration load: ", err)
		return
	}
	if len(args) > 0 {
		runFileStatusFile(args[0], config)
	} else {
		buildProjectStatus(config)
	}
}

func buildProjectStatus(config *general.Config) {
	data := make([][]string, 0, 0)
	response, err := general.GetWorkspace(config)
	if err != nil {
		return
	}
	renderTable2Stdin(headers, data)
}

func runFileStatusFile(fileName string, config *general.Config) {
	data := make([][]string, 0, 0)
	renderTable2Stdin(headers, data)
}
