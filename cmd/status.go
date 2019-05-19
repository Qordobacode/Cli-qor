package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"
	"github.com/spf13/cobra"
	"strconv"
)

var (
	basicHeaders = []string{
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
	workspace, err := general.GetWorkspace(config)
	if err != nil {
		return
	}
	fileSearchResponse := getFileSearchResponse(workspace, config)
	header := buildTableHeader(fileSearchResponse)
	data := buildTableData(fileSearchResponse.ByPersonaProgress, header)
	renderTable2Stdin(header, data)
}

func getFileSearchResponse(response *general.Workspace, config *general.Config) *general.FileSearchResponse {
	for _, person := range response.TargetPersonas {
		fileSearchResponse, err := general.SearchForFiles(config, person.ID, true)
		if err != nil {
			continue
		}
		body, err := json.MarshalIndent(fileSearchResponse, "  ", "  ")
		if err != nil {
			log.Debugf("person %v fileSearchResponse\n%s", person.Code, string(body))
		}
		return fileSearchResponse
	}
	return nil
}

func buildTableHeader(fileSearchResponse *general.FileSearchResponse) []string {
	result := make([]string, 0, len(fileSearchResponse.ByPersonaProgress)+3)
	result = append(result, basicHeaders...)
	for _, personaProgress := range fileSearchResponse.ByPersonaProgress {
		for _, workflowState := range personaProgress.ByWorkflowProgress {
			result = append(result, workflowState.Workflow.Name)
		}
		break
	}
	return result
}

func buildTableData(personaProgress []general.ByPersonaProgress, header []string) [][]string {
	data := make([][]string, 0, len(personaProgress))
	for _, progress := range personaProgress {
		row, err := buildTableRow(&progress, header)
		if err != nil {
			continue
		}
		data = append(data, row)
	}
	return data
}

func buildTableRow(personProgress *general.ByPersonaProgress, header []string) ([]string, error) {
	row := make([]string, len(header))
	row[0] = personProgress.Persona.Code
	totalWords := 0
	totalSegments := 0
	for _, workflowProgress := range personProgress.ByWorkflowProgress {
		totalWords += workflowProgress.Counts.WordCount
		totalSegments += workflowProgress.Counts.SegmentCount
	}
	row[1] = strconv.Itoa(totalWords)
	row[2] = strconv.Itoa(totalSegments)
	i := 0
	// same order in iteration as it was on header filling step
	for _, workflowProgress := range personProgress.ByWorkflowProgress {
		percent := float64(workflowProgress.Counts.SegmentCount) / float64(totalSegments) * 100
		row[i+3] = fmt.Sprintf(`%6.2f%%`, percent)
		i++
	}
	return row, nil
}

func runFileStatusFile(fileName string, config *general.Config) {
	data := make([][]string, 0)
	header := make([]string, 0)
	renderTable2Stdin(header, data)
}
