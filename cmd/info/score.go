package info

import (
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/spf13/cobra"
)

var (
	scoreFileName    string
	scoreFileVersion string
)

func NewScoreCommand() *cobra.Command {
	scoreCommand := &cobra.Command{
		Annotations: map[string]string{"group": "info"},
		Use:         "score",
		Short:       "Score per file",

		PreRun: StartLocalServices,
		Run:    scoreFile,
	}
	scoreCommand.Flags().StringVarP(&scoreFileName, "files", "f", "", "File to score")
	scoreCommand.Flags().StringVarP(&scoreFileVersion, "version", "v", "", "Version of file to score")
	return scoreCommand
}

func scoreFile(cmd *cobra.Command, args []string) {
	if scoreFileName == "" && len(args) > 0 {
		scoreFileName = args[0]
	}
	score := FileService.FileScore(scoreFileName, scoreFileVersion)
	if score != nil {
		log.Infof("%v", score.DocumentScore)
	}
}
