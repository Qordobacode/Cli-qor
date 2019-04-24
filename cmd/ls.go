package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/qordobacode/cli-v2/general"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

const (
	enabled  = "ENABLED"
	disabled = "DISABLED"
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
	data := make([][]string, 0, 0)

	for _, targetPersona := range workspace.TargetPersonas {
		result := handlePersonResult(config, &targetPersona)
		data = append(data, result...)
	}

	render2Stdin(data)
}

func handlePersonResult(config *general.Config, persona *general.Person) [][]string {
	files, e := general.GetFilesForTargetPerson(config, persona.ID)
	data := make([][]string, 0, 0)
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

func buildDataRowFromFile(file *general.File) []string {
	row := make([]string, len(lsHeaders), len(lsHeaders))
	row[0] = strconv.Itoa(file.FileID)
	row[1] = file.Filename
	row[2] = file.Version
	row[3] = strings.Join(file.Tags, ", ")
	row[4] = "" // TODO: add #segments here
	row[5] = general.GetDateFromTimestamp(file.Update)
	row[6] = disabled
	if file.Enabled {
		row[6] = enabled
	}
	return row
}

func render2Stdin(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(lsHeaders)
	table.AppendBulk(data)
	table.Render() // Send output
}
