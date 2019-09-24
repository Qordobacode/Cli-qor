package config

import (
	"bufio"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/config"
	"github.com/qordobacode/cli-v2/pkg/general"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var (
	configurationService pkg.ConfigurationService
)

// NewInitCmd function create `init` command
func NewInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:         "init",
		Short:       "Init configuration for.ConfigConfig CLI from STDIN or file",
		RunE:        RunInitRoot,
		Example:     `"qor init", "qor init qordobaconfig.yaml"`,
		Annotations: map[string]string{"group": "init"},
	}
	var local *general.Local
	configurationService = &config.ConfigurationService{
		Local: local,
	}
	return initCmd
}

// RunInitRoot function starts config initialization
func RunInitRoot(cmd *cobra.Command, args []string) error {
	fmt.Println("init called")
	fileName := ""
	if len(args) > 0 {
		fileName = args[0]
	}
	var newConfig *types.Config
	var err error
	if fileName != "" {
		newConfig, err = configurationService.ReadConfigInPath(fileName)
		if err != nil {
			return err
		}
	}
	if newConfig == nil {
		newConfig = buildConfigFromStdin()
	}

	configurationService.SaveMainConfig(newConfig)
	return nil
}

func buildConfigFromStdin() *types.Config {
	scanner := bufio.NewScanner(os.Stdin)
	accessToken := readVariable("ACCESS TOKEN: ", "Access token can't be empty", scanner)
	organizationID := readIntVariable("ORGANIZATION ID: ", "Organization ID can't be empty", scanner)
	projectID := readIntVariable("WORKSPACE ID: ", "Project ID can't be empty", scanner)
	return &types.Config{
		Qordoba: types.QordobaConfig{
			AccessToken:    accessToken,
			WorkspaceID:    projectID,
			OrganizationID: organizationID,
		},
	}
}

func readVariable(header, errMessage string, scanner *bufio.Scanner) string {
	for {
		fmt.Print(header)
		scanner.Scan()
		text := scanner.Text()
		if text != "" {
			return text
		}
		log.Infof("%s", errMessage)
	}
}

func readIntVariable(header, errMessage string, scanner *bufio.Scanner) int64 {
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
		log.Errorf(errMessage)
	}
}
