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
	workspace, err := general.GetWorkspace(config)
	if err != nil {
		return
	}
	person2FilesMap := getPerson2FilesMap(workspace, config)
	header := buildTableHeader(person2FilesMap)
	data := buildTableData(person2FilesMap, header)

	renderTable2Stdin(header, data)
}

func getPerson2FilesMap(response *general.Workspace, config *general.Config) map[string]*general.FileSearchResponse {
	person2FileMap := make(map[string]*general.FileSearchResponse)
	//var wg sync.WaitGroup
	//wg.Add(len(response.TargetPersonas))
	for _, person := range response.TargetPersonas {
		fileSearchResponse, err := general.SearchForFiles(config, person.ID, true)
		if err != nil {
			continue
		}
		person2FileMap[person.Code] = fileSearchResponse
	}
	//wg.Wait()
	return person2FileMap
}

func buildTableHeader(person2FilesMap map[string]*general.FileSearchResponse) []string {
	workflowName2ColumnMap := make(map[string]int)
	for person, fileSearchResponse := range person2FilesMap {
		for _, byPersonaProgress := range fileSearchResponse.ByPersonaProgress {
			if byPersonaProgress.Persona.Code != person {
				continue
			}
			for _, progress := range byPersonaProgress.ByWorkflowProgress {
				workflowName2ColumnMap[progress.Workflow.Name] = len(workflowName2ColumnMap)
			}
		}
	}
	result := make([]string, 0, len(workflowName2ColumnMap)+3)
	result = append(result, basicHeaders...)
	for key := range workflowName2ColumnMap {
		result = append(result, key)
	}
	return result
}

func buildTableData(person2FilesMap map[string]*general.FileSearchResponse, header []string) [][]string {
	data := make([][]string, 0, len(person2FilesMap))
	for person, fileSearchResponse := range person2FilesMap {
		row, err := buildTableRow(fileSearchResponse, person, header)
		if err != nil {
			continue
		}
		data = append(data, row)
	}
	return data
}

func buildTableRow(fileSearchResponse *general.FileSearchResponse,
	persona string, header []string) ([]string, error) {
	row := make([]string, len(header), len(header))
	row[0] = persona
	row[1] = strconv.Itoa(fileSearchResponse.TotalCounts.WordCount)
	row[2] = strconv.Itoa(fileSearchResponse.TotalCounts.SegmentCount)
	headerMap := make(map[string]int)
	total := 0
	for _, personaProgress := range fileSearchResponse.ByPersonaProgress {
		if personaProgress.Persona.Code != persona {
			continue
		}
		for _, byWorkflowProgress := range personaProgress.ByWorkflowProgress {
			total += byWorkflowProgress.Counts.WordCount
			headerMap[byWorkflowProgress.Workflow.Name] = byWorkflowProgress.Counts.WordCount
		}
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
