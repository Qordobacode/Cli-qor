package segment

import (
	"github.com/spf13/cobra"
)

// deleteKeyCmd represents the deleteKey command
var (
	deleteKeyVersion string
	deleteKeyKey     string
	deleteKeyValue   string
	deleteKeyRef     string
)

func NewDeleteSegmentCommand() *cobra.Command {
	deleteKeyCmd := &cobra.Command{
		Annotations: map[string]string{"group": "segment"},
		Use:         "deleteKey",
		Short:       "Delete segment",
		PreRun:      StartLocalServices,
		Run:         deleteSegment,
	}

	deleteKeyCmd.Flags().StringVarP(&addKeyVersion, "version", "v", "", "file version")
	deleteKeyCmd.Flags().StringVarP(&addKeyKey, "key", "k", "", "key to add")
	deleteKeyCmd.Flags().StringVar(&addKeyValue, "value", "", "value to add")
	deleteKeyCmd.Flags().StringVarP(&addKeyRef, "ref", "r", "", "")
	return deleteKeyCmd
}

func deleteSegment(cmd *cobra.Command, args []string) {

}
