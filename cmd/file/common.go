package file

import (
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/config"
	"github.com/qordobacode/cli-v2/pkg/file"
	"github.com/qordobacode/cli-v2/pkg/general"
	"github.com/qordobacode/cli-v2/pkg/rest"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/qordobacode/cli-v2/pkg/workspace"
	"github.com/spf13/cobra"
	"os"
)

var (
	appConfig            *types.Config
	local                *general.Local
	configurationService = config.ConfigurationService{
		Local: local,
	}
	qordobaClient    pkg.QordobaClient
	workspaceService pkg.WorkspaceService
	fileService      pkg.FileService
)

// startLocalServices function build all required for file package services
func startLocalServices(cmd *cobra.Command, args []string) {
	var err error
	appConfig, err = configurationService.LoadConfig()
	if err != nil {
		os.Exit(1)
	}
	qordobaClient = rest.NewRestClient(appConfig)
	workspaceService = &workspace.Service{
		Config:        appConfig,
		QordobaClient: qordobaClient,
		Local:         local,
	}
	fileService = &file.Service{
		Config:           appConfig,
		WorkspaceService: workspaceService,
		Local:            local,
		QordobaClient:    qordobaClient,
	}
	local = &general.Local{
		Config: appConfig,
	}
}
