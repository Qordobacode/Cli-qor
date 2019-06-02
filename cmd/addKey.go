package cmd

import (
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"

	"github.com/spf13/cobra"
)

var (
	// addKeyCmd represents the add-key command
	addKeyCmd = &cobra.Command{
		Use:     "add-key",
		Short:   "A brief description of your command",
		PreRunE: preValidateParameters,
		Run:     addKey,
	}
	addKeyVersion string
	addKeyKey     string
	addKeyValue   string
	addKeyRef     string
)

func preValidateParameters(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("filename is mandatory")
	}
	if addKeyKey == "" {
		return fmt.Errorf("flag 'key' is mandatory")
	}
	if addKeyValue == "" {
		return fmt.Errorf("flag 'value' is mandatory")
	}
	startLocalServices(cmd, args)
	return nil
}

func addKey(cmd *cobra.Command, args []string) {
	log.Debugf("addKey called")
	if Config == nil {
		log.Errorf("error occurred on configuration load")
		return
	}
	keyAddRequest := &types.KeyAddRequest{
		Key:       addKeyKey,
		Source:    addKeyValue,
		Reference: addKeyRef,
	}
	SegmentService.AddKey(args[0], addKeyVersion, keyAddRequest)
}

func init() {
	rootCmd.AddCommand(addKeyCmd)

	addKeyCmd.Flags().StringVarP(&addKeyVersion, "version", "v", "", "file version")
	addKeyCmd.Flags().StringVarP(&addKeyKey, "key", "k", "", "key to add")
	addKeyCmd.Flags().StringVar(&addKeyValue, "value", "", "value to add")
	addKeyCmd.Flags().StringVarP(&addKeyRef, "ref", "r", "", "")
}
