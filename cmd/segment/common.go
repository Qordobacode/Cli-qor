package segment

import (
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/config"
	"github.com/qordobacode/cli-v2/pkg/file"
	"github.com/qordobacode/cli-v2/pkg/general"
	"github.com/qordobacode/cli-v2/pkg/rest"
	"github.com/qordobacode/cli-v2/pkg/segments"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/qordobacode/cli-v2/pkg/workspace"
	"github.com/spf13/cobra"
)

var (
	Config               *types.Config
	Local                *general.Local
	ConfigurationService = config.ConfigurationService{
		Local: Local,
	}
	QordobaClient    pkg.QordobaClient
	WorkspaceService pkg.WorkspaceService
	FileService      pkg.FileService
	SegmentService   pkg.SegmentService
)

func StartLocalServices(cmd *cobra.Command, args []string) {
	var err error
	Config, err = ConfigurationService.LoadConfig()
	if err != nil {
		return
	}
	QordobaClient = rest.NewRestClient(Config)
	WorkspaceService = &workspace.WorkspaceService{
		Config:        Config,
		QordobaClient: QordobaClient,
		Local:         Local,
	}
	FileService = &file.FileService{
		Config:           Config,
		WorkspaceService: WorkspaceService,
		Local:            Local,
		QordobaClient:    QordobaClient,
	}
	SegmentService = &segments.SegmentService{
		QordobaClient:    QordobaClient,
		WorkspaceService: WorkspaceService,
		Config:           Config,
		FileService:      FileService,
	}
}
