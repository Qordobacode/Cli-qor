package cmd

import (
	"bufio"
	"fmt"
	"github.com/qordobacode/cli-v2/general"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:         "init",
	Short:       "Init configuration for QordobaConfig CLI from STDIN",
	RunE:        RunInitRoot,
	Example:     "qor init",
	Annotations: map[string]string{"version": APIVersion},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// RunInitRoot function starts config initialization
func RunInitRoot(cmd *cobra.Command, args []string) error {
	fmt.Println("init called")
	fileName := ""
	if len(args) > 0 {
		fileName = args[0]
	}
	return RunInit(fileName)
}

// RunInit creates a config with the given options.
func RunInit(fileName string) error {
	var config *general.Config
	var err error
	if fileName != "" {
		config, err = general.ReadConfigInPath(fileName)
		if err != nil {
			return err
		}
	}
	if config == nil {
		scanner := bufio.NewScanner(os.Stdin)
		accessToken := readVariable("ACCESS TOKEN: ", "Access token can't be empty\n", scanner)
		organizationID := readIntVariable("ORGANIZATION ID: ", "Organization ID can't be empty\n", scanner)
		projectID := readIntVariable("PROJECT ID: ", "Project ID can't be empty\n", scanner)
		config = &general.Config{
			Qordoba: general.QordobaConfig{
				AccessToken:    accessToken,
				ProjectID:      projectID,
				OrganizationID: organizationID,
			},
		}
	}

	general.SaveMainConfig(config)
	return nil
}

func readVariable(header, errMesage string, scanner *bufio.Scanner) string {
	for {
		fmt.Print(header)
		scanner.Scan()
		text := scanner.Text()
		if text != "" {
			return text
		}
		fmt.Printf(errMesage)
	}
}

func readIntVariable(header, errMesage string, scanner *bufio.Scanner) int64 {
	for {
		fmt.Print(header)
		scanner.Scan()
		text := scanner.Text()
		if text != "" {
			num, err := strconv.ParseInt(text, 10, 64)
			if err == nil {
				return num
			}
		}
		fmt.Printf(errMesage)
	}
}
