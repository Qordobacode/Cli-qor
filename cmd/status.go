package cmd

import (
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"
	"github.com/spf13/cobra"
)

var (
	statusHeaders = []string{
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
	response, err := general.GetWorkspace(config)
	if err != nil {
		return
	}
	data := make([][]string, 0, len(response.TargetPersonas))
	for _, person := range response.TargetPersonas {
		row := make([]string, len(statusHeaders), len(statusHeaders))
		row[0] = person.Code
		personFiles, err := general.GetFilesForTargetPerson(config, person.ID, true)
		if err != nil {
			continue
		}

		data = append(data, row)
	}

	renderTable2Stdin(statusHeaders, data)
}

func runFileStatusFile(fileName string, config *general.Config) {
	data := make([][]string, 0, 0)
	renderTable2Stdin(statusHeaders, data)
}
