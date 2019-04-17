package cmd

import (
	"bufio"
	"fmt"
	"github.com/qordobacode/cli-v2/files"
	"github.com/qordobacode/cli-v2/models"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init configuration for Qordoba CLI from STDIN",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("init called")
		fileName := ""
		if len(args) > 0 {
			fileName = args[0]
		}
		return RunInit(fileName)
	},
	Example:     "qor init",
	Annotations: map[string]string{"version": ApiVersion},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// RunConfigCreate creates a config with the given options.
func RunInit(fileName string) error {
	var config *models.QordobaConfig
	var err error
	if fileName != "" {
		config, err = files.ReadConfigInPath(fileName)
		if err != nil {
			return err
		}
	}
	if config == nil {
		scanner := bufio.NewScanner(os.Stdin)
		accessToken := readVariable("ACCESS TOKEN: ", "Access token can't be empty\n", scanner)
		organizationID := readIntVariable("ORGANIZATION ID: ", "Organization ID can't be empty\n", scanner)
		projectID := readIntVariable("PROJECT ID: ", "Project ID can't be empty\n", scanner)
		config = &models.QordobaConfig{
			Qordoba: models.Qordoba{
				AccessToken:    accessToken,
				ProductID:      projectID,
				OrganizationID: organizationID,
			},
		}
	}

	files.SaveMainConfig(config)
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
