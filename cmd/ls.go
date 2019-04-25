package cmd

import (
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	enabled   = "ENABLED"
	disabled  = "DISABLED"
	lineLimit = 50
)

// lsCmd represents the ls command
var (
	lsCmd = &cobra.Command{
		Use:   "ls",
		Short: "Ls files (show 50 only)",
		Run:   printLs,
	}
	lsHeaders = []string{"ID", "NAME", "version", "tag", "#SEGMENTS", "UPDATED_ON", "STATUS"}
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

func printLs(cmd *cobra.Command, args []string) {
	config, e := general.LoadConfig()
	if e != nil {
		return
	}
	workspace, err := general.GetWorkspace(config)
	if err != nil {
		return
	}
	data := make([]*responseRow, 0, 0)
	for _, targetPersona := range workspace.TargetPersonas {
		result := handlePersonResult(config, &targetPersona)
		data = append(data, result...)
		if len(data) > 50 {
			data = data[:50]
			break
		}
	}
	// add sorting for output
	sort.Slice(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
	})

	printFile2Stdin(data)
}

func printFile2Stdin(data []*responseRow) {
	if !IsJSON {
		render2Stdin(data)
	} else {
		bytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Errorf("error occurred on marshalling with JSON: %v", err)
			return
		}
		log.Infof("%v", string(bytes))
	}
}

func handlePersonResult(config *general.Config, persona *general.Person) []*responseRow {
	files, e := general.GetFilesForTargetPerson(config, persona.ID)
	data := make([]*responseRow, 0, 0)
	if e != nil {
		return data
	}
	audiences := config.GetAudiences()
	for _, file := range files {
		if _, ok := audiences[persona.Code]; len(audiences) > 0 && !ok {
			continue
		}
		data = append(data, buildDataRowFromFile(&file))
	}
	return data
}

func buildDataRowFromFile(file *general.File) *responseRow {
	// strconv.Itoa
	row := responseRow{
		ID:        file.FileID,
		Name:      file.Filename,
		Version:   file.Version,
		Tag:       file.Tags,
		UpdatedOn: general.GetDateFromTimestamp(file.Update),
		Status:    disabled,
	}
	if file.Enabled {
		row.Status = enabled
	}
	return &row
}

func render2Stdin(response []*responseRow) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(lsHeaders)
	data := formatResponse2Array(response)
	table.AppendBulk(data)
	table.Render() // Send output
}

func formatResponse2Array(rows []*responseRow) [][]string {
	data := make([][]string, 0, len(rows))
	for _, responseRow := range rows {
		row := make([]string, len(lsHeaders), len(lsHeaders))
		row[0] = strconv.Itoa(responseRow.ID)
		row[1] = responseRow.Name
		row[2] = responseRow.Version
		row[3] = strings.Join(responseRow.Tag, ", ")
		row[4] = "" // TODO: add #segments here
		row[5] = responseRow.UpdatedOn
		row[6] = responseRow.Status
		data = append(data, row)
	}
	return data
}

// responseRow struct
type responseRow struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Tag         []string `json:"tag"`
	SegmentNums int      `json:"#segments"`
	UpdatedOn   string   `json:"updated_on"`
	Status      string   `json:"status"`
}
