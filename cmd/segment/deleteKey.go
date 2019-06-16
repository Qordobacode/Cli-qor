package segment

import (
	"github.com/spf13/cobra"
)

// deleteKeyCmd represents the deleteKey command
var (
	deleteKeyVersion string
	deleteKeyKey     string
)

// NewDeleteSegmentCommand function add `delete-key` command
func NewDeleteSegmentCommand() *cobra.Command {
	deleteKeyCmd := &cobra.Command{
		Annotations: map[string]string{"group": "segment"},
		Use:         "delete-key",
		Short:       "Delete segment",
		Example:     `qor delete-key file_name.doc --version v1 --key "/go_nav_menu"`,
		PreRun:      startLocalServices,
		Run:         deleteSegment,
	}

	deleteKeyCmd.Flags().StringVarP(&deleteKeyVersion, "version", "v", "", "file version where update segment")
	deleteKeyCmd.Flags().StringVarP(&deleteKeyKey, "key", "k", "", "key to delete")
	return deleteKeyCmd
}

func deleteSegment(cmd *cobra.Command, args []string) {
	segmentService.DeleteKey(args[0], deleteKeyVersion, deleteKeyKey)
}
