package cmd

import (
	"fmt"
	"github.com/qordobacode/cli-v2/general"
	"github.com/qordobacode/cli-v2/log"

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
	return nil
}

func addKey(cmd *cobra.Command, args []string) {
	config, e := general.LoadConfig()
	if e != nil {
		return
	}
	log.Debugf("addKey called")
	keyAddRequest := &general.KeyAddRequest{
		Key:       addKeyKey,
		Source:    addKeyValue,
		Reference: addKeyRef,
	}
	general.AddKey(config, args[0], addKeyVersion, keyAddRequest)
}

func init() {
	rootCmd.AddCommand(addKeyCmd)

	addKeyCmd.Flags().StringVarP(&addKeyVersion, "version", "v", "", "file version")
	addKeyCmd.Flags().StringVarP(&addKeyKey, "key", "k", "", "key to add")
	addKeyCmd.Flags().StringVar(&addKeyValue, "value", "", "value to add")
	addKeyCmd.Flags().StringVarP(&addKeyRef, "ref", "r", "", "")
}
