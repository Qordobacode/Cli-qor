package segment

import (
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/spf13/cobra"
	"strconv"
)

// updateValueCmd represents the updateValue command
var (
	valueKeyVersion string
	valueKeyKey     string
	IsJSON          bool
	header          = []string{"FILE NAME", "#VERSION", "KEY", "#VALUE", "#REF", "#TIMESTAMP"}
)

func NewValueKeyCommand() *cobra.Command {
	valueKeyCmd := &cobra.Command{
		Annotations: map[string]string{"group": "segment"},
		Use:         `value-key`,
		Example:     `qor value-key file_name.doc --version v1 --key "/go_nav_menu"`,
		Short:       "Pull value by key",
		PreRunE:     preValidateValueKeyParameters,
		Run:         pullValueByKey,
	}

	valueKeyCmd.Flags().StringVarP(&valueKeyVersion, "version", "v", "", "file version")
	valueKeyCmd.Flags().StringVarP(&valueKeyKey, "key", "k", "", "key to get value")
	valueKeyCmd.PersistentFlags().BoolVar(&IsJSON, "json", false, "Print output in JSON format")
	return valueKeyCmd
}

func preValidateValueKeyParameters(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("filename is mandatory")
	}
	if valueKeyKey == "" {
		return fmt.Errorf("flag 'key' is mandatory")
	}
	StartLocalServices(cmd, args)
	return nil
}
func pullValueByKey(cmd *cobra.Command, args []string) {
	segment, _ := SegmentService.FindSegment(args[0], valueKeyVersion, valueKeyKey)
	resultArray := make([]*valueInfo, 0)
	resultArray = append(resultArray, &valueInfo{
		Filename:    args[0],
		FileVersion: valueKeyVersion,
		Key:         valueKeyKey,
		Value:       segment.SsText,
		Reference:   segment.Reference,
		Timestamp:   strconv.Itoa(segment.LastSaved),
	})
	printProjectStatus2Stdin(resultArray)
}

type valueInfo struct {
	Filename    string `json:"file_name"`
	FileVersion string `json:"file_version"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Reference   string `json:"reference"`
	Timestamp   string `json:"timestamp"`
}

func printProjectStatus2Stdin(valueInfo []*valueInfo) {
	if !IsJSON {
		data := formatResponse2Array(valueInfo)
		Local.RenderTable2Stdin(header, data)
		return
	}
	bytes, err := json.MarshalIndent(valueInfo, "", "  ")
	if err != nil {
		log.Errorf("error occurred on marshalling with JSON: %v", err)
		return
	}
	log.Infof("%v", string(bytes))
}

func formatResponse2Array(response []*valueInfo) [][]string {
	data := make([][]string, 0, len(response))
	for _, valInfo := range response {
		row := make([]string, len(header))
		header = []string{"FILE NAME", "#VERSION", "KEY", "#VALUE", "#REF", "#TIMESTAMP"}
		row[0] = valInfo.Filename
		row[1] = valInfo.FileVersion
		row[2] = valInfo.Key
		row[3] = valInfo.Value
		row[4] = valInfo.Reference
		row[5] = valInfo.Timestamp
	}
	return data
}
