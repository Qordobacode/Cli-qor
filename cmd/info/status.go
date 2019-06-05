package info

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

func NewStatusCommand() *cobra.Command{
	statusCmd := &cobra.Command{
		Annotations: map[string]string{"group": "info"},
		Use:   "status",
		Short: "Status per project or file (Support file versions)",
		Run:   runStatus,
		PreRun: StartLocalServices,
	}
	statusCmd.Flags().StringVarP(&statusFileVersion, "version", "v", "", "--version")
	statusCmd.PersistentFlags().BoolVar(&IsJSON, "json", false, "Print output in JSON format")
	return statusCmd

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
	if fileSearchResponse == nil {
		log.Errorf("fileSearchResponse %v not found", workspace.Workspace.ID)
		return
	}
	if len(fileSearchResponse.ByPersonaProgress) == 0 {
		log.Errorf("fileSearchResponse returned empty byPersonaProgress. Status is not available")
		return
	}
	progress := fileSearchResponse.ByPersonaProgress[0].ByWorkflowProgress
	header := buildTableHeader(progress, basicHeaders)
	dataRow, dataJson := buildTableData(fileSearchResponse.ByPersonaProgress, header)
	printProjectStatus2Stdin(header, dataRow, dataJson)
}

func printProjectStatus2Stdin(header []string, tableData [][]string, dataJson []map[string]string) {
	if !IsJSON {
		renderTable2Stdin(header, tableData)
	} else {
		bytes, err := json.MarshalIndent(dataJson, "", "  ")
		if err != nil {
			log.Errorf("error occurred on marshalling with JSON: %v", err)
			return
		}
		log.Infof("%v", string(bytes))
	}
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

func buildTableData(personaProgress []types.ByPersonaProgress, header []string) ([][]string, []map[string]string) {
	data := make([][]string, 0, len(personaProgress))
	dataJSON := make([]map[string]string, 0, len(personaProgress))
	for _, progress := range personaProgress {
		row, rowJSON := buildTableRow(&progress, header)
		data = append(data, row)
		dataJSON = append(dataJSON, rowJSON)
	}
	return data, dataJSON
}

func buildTableRow(personProgress *types.ByPersonaProgress, header []string) ([]string, map[string]string) {
	row := make([]string, len(header))
	rowMap := make(map[string]string)
	row[0] = personProgress.Persona.Code
	rowMap["audiences_id"] = strconv.Itoa(personProgress.Persona.ID)
	rowMap["audiences_code"] = personProgress.Persona.Code
	rowMap["audiences_name"] = personProgress.Persona.Name

	totalWords := 0
	totalSegments := 0
	for _, workflowProgress := range personProgress.ByWorkflowProgress {
		totalWords += workflowProgress.Counts.WordCount
		totalSegments += workflowProgress.Counts.SegmentCount
	}
	row[1] = strconv.Itoa(totalWords)
	rowMap["words_total"] = row[1]
	row[2] = strconv.Itoa(totalSegments)
	rowMap["segments_total"] = row[2]
	i := 0
	// same order in iteration as it was on header filling step
	for _, workflowProgress := range personProgress.ByWorkflowProgress {
		percent := float64(workflowProgress.Counts.SegmentCount) / float64(totalSegments) * 100
		row[i+3] = fmt.Sprintf(`%6.2f%%`, percent)
		rowMap[workflowProgress.Workflow.Name] = row[i+3]
		i++
	}
	return row, rowMap
}

func runFileStatusFile(fileName string, workspace *types.WorkspaceData) {
	fileSearchResponse, _ := FileService.FindFile(fileName, statusFileVersion, true)
	if fileSearchResponse == nil {
		log.Debugf("file %s('%s') was not found", fileName, statusFileVersion)
		return
	}
	header := buildTableHeader(fileSearchResponse.ByWorkflowProgress, fileHeaders)
	data, dataJSON := buildFileTableData(fileSearchResponse)
	printProjectStatus2Stdin(header, data, dataJSON)
}

func buildFileTableData(file *types.File) ([][]string, []map[string]string) {
	data := make([][]string, 0, len(file.ByWorkflowProgress))
	dataJSON := make([]map[string]string, 0, len(file.ByWorkflowProgress))
	fileRow := make([]string, 0, 0)
	JSONdocument := make(map[string]string)
	fileRow[0] = file.Filename + " " + file.Version
	fileRow[1] = strconv.Itoa(file.Counts.WordCount)
	JSONdocument["words"] = fileRow[1]
	fileRow[2] = strconv.Itoa(file.Counts.SegmentCount)
	JSONdocument["segments"] = fileRow[2]
	JSONdocument["file_name"] = file.Filename
	JSONdocument["version"] = file.Version
	JSONdocument["file_id"] = strconv.Itoa(file.FileID)
	if file.Completed {
		JSONdocument["completed"] = strconv.FormatBool(file.Completed)
	}
	i := 0
	for _, workflowProgress := range file.ByWorkflowProgress {
		percent := float64(workflowProgress.Counts.SegmentCount) / float64(file.Counts.SegmentCount) * 100
		fileRow[i+3] = fmt.Sprintf(`%6.2f%%`, percent)
		JSONdocument[workflowProgress.Workflow.Name] = fileRow[i+3]
		i++
	}

	return data, dataJSON
}
