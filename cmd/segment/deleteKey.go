package segment

import (
	"github.com/spf13/cobra"
)

// deleteKeyCmd represents the deleteKey command
var (
	deleteKeyVersion string
	deleteKeyKey     string
)

func NewDeleteSegmentCommand() *cobra.Command {
	deleteKeyCmd := &cobra.Command{
		Annotations: map[string]string{"group": "segment"},
		Use:         "delete-key",
		Short:       "Delete segment",
		Example:     `qor delete-key file_name.doc --version v1 --key "/go_nav_menu"`,
		PreRun:      StartLocalServices,
		Run:         deleteSegment,
	}

	deleteKeyCmd.Flags().StringVarP(&deleteKeyVersion, "version", "v", "", "file version where update segment")
	deleteKeyCmd.Flags().StringVarP(&deleteKeyKey, "key", "k", "", "key to delete")
	return deleteKeyCmd
}

func deleteSegment(cmd *cobra.Command, args []string) {
	SegmentService.DeleteKey(args[0], deleteKeyVersion, deleteKeyKey)
}
