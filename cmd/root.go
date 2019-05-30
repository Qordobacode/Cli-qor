package cmd

import (
	"fmt"
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/config"
	"github.com/qordobacode/cli-v2/pkg/file"
	"github.com/qordobacode/cli-v2/pkg/general"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/rest"
	"github.com/qordobacode/cli-v2/pkg/segments"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/qordobacode/cli-v2/pkg/workspace"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:     "qor",
		Short:   "Qordoba CLI",
		Long:    `This CLI is used for simplified access to Qordoba API`,
		Version: APIVersion,
		Run: func(cmd *cobra.Command, args []string) {
			if Version {
				printVersion()
			}
		},
	}
	Help                 bool
	Version              bool
	IsJSON               bool
	Local                *general.Local
	ConfigurationService = config.ConfigurationService{
		Local: Local,
	}
	Config           *types.Config
	QordobaClient    pkg.QordobaClient
	WorkspaceService pkg.WorkspaceService
	FileService      pkg.FileService
	SegmentService   pkg.SegmentService
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Version, "version", "v", false, "Get version of CLI")
	rootCmd.PersistentFlags().BoolVar(&log.IsVerbose, "verbose", false, "Print verbose output")
	rootCmd.PersistentFlags().BoolVar(&IsJSON, "json", false, "Print output in JSON format")
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

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".test" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".test")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
