package file

import (
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
)

var (
	scoreGetTemplate = "%s/v3/contentscore/organizations/%d/workspaces/%d/documents/%v/personas/%d/score?documentLength=%d"
)

// FileScore function returns file score
func (f *Service) FileScore(filename, version string) *types.ScoreResponseBody {
	file, personaID := f.FindFile(filename, version, false)
	if file == nil {
		return nil
	}
	base := f.Config.GetAPIBase()
	fileListURL := fmt.Sprintf(scoreGetTemplate, base, f.Config.Qordoba.OrganizationID,
		f.Config.Qordoba.WorkspaceID, file.FileID, personaID, 1)
	sourceResponse, err := f.QordobaClient.GetFromServer(fileListURL)
	if err != nil {
		return nil
	}
	var scoreResponseBody types.ScoreResponseBody
	err = scoreResponseBody.UnmarshalJSON(sourceResponse)
	if err != nil {
		log.Errorf("error occurred on file score response unmarshalling: %v", err)
		return nil
	}
	return &scoreResponseBody
}
