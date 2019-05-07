package cmd

import (
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
	response, err := general.GetWorkspace(config)
	if err != nil {
		return
	}
	header := buildTableHeader(response, config)
	data := buildTableData(response, config, header)

	renderTable2Stdin(header, data)
}

func buildTableHeader(response *general.Workspace, config *general.Config) []string {
	workflowName2ColumnMap := make(map[string]int)
	for _, person := range response.TargetPersonas {
		fileSearchResponse, err := general.SearchForFiles(config, person.ID, true)
		if err != nil {
			continue
		}
		personaProgress := getPersonaProgress(fileSearchResponse, person)
		if personaProgress == nil {
			continue
		}
		updateWorkflowColumnMap(personaProgress, workflowName2ColumnMap)
	}
	result := make([]string, 0, len(workflowName2ColumnMap) + 3)
	result = append(result, basicHeaders...)
	for key := range workflowName2ColumnMap {
		result = append(result, key)
	}
	return result
}

func buildTableData(response *general.Workspace, config *general.Config, header []string) [][]string {
	data := make([][]string, 0, len(response.TargetPersonas))
	for _, person := range response.TargetPersonas {
		fileSearchResponse, err := general.SearchForFiles(config, person.ID, true)
		if err != nil {
			continue
		}
		personaProgress := getPersonaProgress(fileSearchResponse, person)
		if personaProgress == nil {
			continue
		}
		row, err := buildTableRow(&person, fileSearchResponse, personaProgress, header)
		if err != nil {
			continue
		}
		data = append(data, row)
	}
	return data
}

func getPersonaProgress(fileSearchResponse *general.FileSearchResponse, person general.Person) *general.ByPersonaProgress {
	for _, personaElem := range fileSearchResponse.ByPersonaProgress {
		if personaElem.Persona.ID == person.ID {
			return &personaElem
		}
	}
	return nil
}

func updateWorkflowColumnMap(personaProgress *general.ByPersonaProgress, workflowName2ColumnMap map[string]int) {
	for _, byWorkflowProgress := range personaProgress.ByWorkflowProgress {
		if _, ok := workflowName2ColumnMap[byWorkflowProgress.Workflow.Name]; !ok {
			workflowName2ColumnMap[byWorkflowProgress.Workflow.Name] = len(workflowName2ColumnMap)
		}
	}
}

func buildTableRow(person *general.Person, fileSearchResponse *general.FileSearchResponse,
	personProgress *general.ByPersonaProgress, header []string) ([]string, error) {
	row := make([]string, len(header), len(header))
	row[0] = person.Code
	row[1] = strconv.Itoa(fileSearchResponse.TotalCounts.WordCount)
	row[2] = strconv.Itoa(fileSearchResponse.TotalCounts.SegmentCount)
	headerMap := make(map[string]int)
	total := 0
	for _, byWorkflowProgress := range personProgress.ByWorkflowProgress {
		total += byWorkflowProgress.Counts.WordCount
		headerMap[byWorkflowProgress.Workflow.Name] = byWorkflowProgress.Counts.WordCount
	}
	for i := 3; i < len(header); i++ {
		statusCount := headerMap[header[i]]
		percent := float64(statusCount) / float64(total) * 100
		row[i] = fmt.Sprintf("%v%", percent)
	}

	return row, nil
}

func runFileStatusFile(fileName string, config *general.Config) {
	data := make([][]string, 0, 0)
	header := make([]string, 0, 0)
	renderTable2Stdin(header, data)
}
