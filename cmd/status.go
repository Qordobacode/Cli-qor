package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"

	"github.com/spf13/cobra"
	"strconv"
)

var (
	basicHeaders = []string{
		"#AUDIENSES", "#WORDS", "#SEGMENTS",
	}
	fileHeaders = []string{
		"FILE NAME", "#WORDS", "#SEGMENTS",
	}
	statusFileVersion string
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status per project or file (Support file versions)",
	Run:   runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	pushCmd.Flags().StringVarP(&statusFileVersion, "version", "v", "", "--version")
}

func runStatus(cmd *cobra.Command, args []string) {
	if Config == nil {
		log.Errorf("error occurred on configuration load")
		return
	}
	workspace, err := WorkspaceService.LoadWorkspace()
	if err != nil || workspace == nil {
		return
	}
	if len(args) > 0 {
		runFileStatusFile(args[0], workspace)
	} else {
		buildProjectStatus(workspace)
	}
}

func buildProjectStatus(workspace *types.WorkspaceData) {
	fileSearchResponse := getFileSearchResponse(&workspace.Workspace)
	progress := fileSearchResponse.ByPersonaProgress[0].ByWorkflowProgress
	header := buildTableHeader(progress, basicHeaders)
	data := buildTableData(fileSearchResponse.ByPersonaProgress, header)
	renderTable2Stdin(header, data)
}

func getFileSearchResponse(response *types.Workspace) *types.FileSearchResponse {
	for _, person := range response.TargetPersonas {
		fileSearchResponse, err := FileService.WorkspaceFiles(person.ID, true)
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

func buildTableHeader(workflowProgress []types.ByWorkflowProgress, headers []string) []string {
	result := make([]string, 0, len(workflowProgress)+len(headers))
	result = append(result, headers...)
	for _, workflowState := range workflowProgress {
		result = append(result, workflowState.Workflow.Name)
	}
	return result
}

func buildTableData(personaProgress []types.ByPersonaProgress, header []string) [][]string {
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

func buildTableRow(personProgress *types.ByPersonaProgress, header []string) ([]string, error) {
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

func runFileStatusFile(fileName string, workspace *types.WorkspaceData) {
	fileSearchResponse, _ := FileService.FindFile(fileName, statusFileVersion, true)
	if fileSearchResponse == nil {
		log.Debugf("file %s('%s') was not found", fileName, statusFileVersion)
		return
	}
	header := buildTableHeader(fileSearchResponse.ByWorkflowProgress, fileHeaders)
	data := make([][]string, 0)
	renderTable2Stdin(header, data)
}
