package cmd

import (
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"
	"github.com/spf13/cobra"
	"strconv"
)

var (
	statusHeaders = []string{
		"#AUDIENSES", "#WORDS", "#SEGMENTS",
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
	for i, person := range response.TargetPersonas {
		fileSearchResponse, err := general.SearchForFiles(config, person.ID, true)
		if err != nil {
			continue
		}
		var personaProgress *general.ByPersonaProgress
		for _, personaElem := range fileSearchResponse.ByPersonaProgress {
			if personaElem.Persona.ID == person.ID {
				personaProgress = &personaElem
				break
			}
		}
		if personaProgress == nil {
			continue
		}
		statusHeaders = updateTableHeader(i, personaProgress)
		row, err := buildTableRow(&person, fileSearchResponse, personaProgress)
		if err != nil {
			continue
		}
		data = append(data, row)
	}

	renderTable2Stdin(statusHeaders, data)
}

func updateTableHeader(i int, personaProgress *general.ByPersonaProgress) []string {
	if i == 0 {
		for _, byWorkflowProgress := range personaProgress.ByWorkflowProgress {
			statusHeaders = append(statusHeaders, byWorkflowProgress.Workflow.Name)
		}
	}
	return statusHeaders
}

func buildTableRow(person *general.Person, fileSearchResponse *general.FileSearchResponse, personProgress *general.ByPersonaProgress) ([]string, error) {
	row := make([]string, len(statusHeaders), len(statusHeaders))
	row[0] = person.Code
	row[1] = strconv.Itoa(fileSearchResponse.TotalCounts.WordCount)
	row[2] = strconv.Itoa(fileSearchResponse.TotalCounts.SegmentCount)
	headerMap := make(map[string]int)
	total := 0
	for _, byWorkflowProgress := range personProgress.ByWorkflowProgress {
		total += byWorkflowProgress.Counts.WordCount
		headerMap[byWorkflowProgress.Workflow.Name] = byWorkflowProgress.Counts.WordCount
	}

	return row, nil
}

func runFileStatusFile(fileName string, config *general.Config) {
	data := make([][]string, 0, 0)
	renderTable2Stdin(statusHeaders, data)
}
