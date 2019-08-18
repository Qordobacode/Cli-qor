package info

import (
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"sort"

	"github.com/spf13/cobra"
	"strconv"
)

const (
	completed = "Complete"
)

var (
	basicHeaders = []string{
		"#AUDIENSES", "#SEGMENTS", "#WORDS",
	}
	fileHeaders = []string{
		"FILE NAME", "#SEGMENTS", "#WORDS",
	}
	statusFileVersion string
)

// NewStatusCommand creates `status` command
func NewStatusCommand() *cobra.Command {
	statusCmd := &cobra.Command{
		Annotations: map[string]string{"group": "info"},
		Use:         "status",
		Short:       "Status per project or file (Support file versions)",
		Example:     `"qor status", "qor status --json", "qor status filename.docx --version 0.2"`,
		Run:         runStatus,
		PreRun:      startLocalServices,
	}
	statusCmd.Flags().StringVarP(&statusFileVersion, "version", "v", "", "--version")
	statusCmd.PersistentFlags().BoolVar(&IsJSON, "json", false, "Print output in JSON format")
	return statusCmd

}

func runStatus(cmd *cobra.Command, args []string) {
	if appConfig == nil {
		log.Errorf("error occurred on configuration load")
		return
	}
	workspace, err := workspaceService.LoadWorkspace()
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
	reorderProgressStatuses(fileSearchResponse)
	header := buildTableHeader(progress, basicHeaders)
	dataRow, dataJSON := buildTableData(fileSearchResponse.ByPersonaProgress, header)
	printProjectStatus2Stdin(header, dataRow, dataJSON)
}

func reorderProgressStatuses(fileSearchResponse *types.FileSearchResponse) {
	progresses := fileSearchResponse.ByPersonaProgress
	for k, _ := range progresses {
		sort.Slice(progresses[k].ByWorkflowProgress, func(i, j int) bool {
			return progresses[k].ByWorkflowProgress[i].Workflow.Order < progresses[k].ByWorkflowProgress[j].Workflow.Order
		})
	}
}

func printProjectStatus2Stdin(header []string, tableData [][]string, dataJSON []map[string]string) {
	if !IsJSON {
		local.RenderTable2Stdin(header, tableData)
		return
	}
	bytes, err := json.MarshalIndent(dataJSON, "", "  ")
	if err != nil {
		log.Errorf("error occurred on marshalling with JSON: %v", err)
		return
	}
	log.Infof("%v", string(bytes))
}

func getFileSearchResponse(response *types.Workspace) *types.FileSearchResponse {
	for _, person := range response.TargetPersonas {
		fileSearchResponse, err := fileService.WorkspaceFiles(person.ID, true)
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
	completedVal := false
	for _, workflowState := range workflowProgress {
		if workflowState.Workflow.Complete {
			completedVal = true
		} else {
			result = append(result, workflowState.Workflow.Name)
		}
	}
	if completedVal {
		result = append(result, completed)
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
	row[1] = strconv.Itoa(totalSegments)
	rowMap["segments_total"] = row[1]
	row[2] = strconv.Itoa(totalWords)
	rowMap["words_total"] = row[2]
	i := 0
	completedVal := ""
	// same order in iteration as it was on header filling step
	for _, workflowProgress := range personProgress.ByWorkflowProgress {
		percent := float64(workflowProgress.Counts.SegmentCount) / float64(totalSegments) * 100
		val := fmt.Sprintf(`%6.2f%%`, percent)
		if workflowProgress.Workflow.Complete {
			///////rowMap[completed] = val*
			completedVal = val
		} else {
			row[i+3] = val
			rowMap[workflowProgress.Workflow.Name] = row[i+3]
			i++
		}
	}
	if completedVal != "" {
		row[i+3] = completedVal
	}
	return row, rowMap
}

func runFileStatusFile(fileName string, workspace *types.WorkspaceData) {
	fileSearchResponse, _ := fileService.FindFile(fileName, statusFileVersion, true)
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
	fileRow[1] = strconv.Itoa(file.Counts.SegmentCount)
	JSONdocument["segments"] = fileRow[1]
	fileRow[2] = strconv.Itoa(file.Counts.WordCount)
	JSONdocument["words"] = fileRow[2]
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
