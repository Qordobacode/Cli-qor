package segments

import (
	"encoding/json"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"io/ioutil"
	"net/http"
)

var (
	keyAddTemplate     = "%s/v3/organizations/%d/workspaces/%d/files/%d/segments/keyAdd"
	getSegmentTemplate = "%s/v3/organizations/%d/workspaces/%d/personas/%v/files/%d/workflow/%d/segments?search=%s"
	keyUpdateTemplate  = "%s/v3/organizations/%d/workspaces/%d/files/%d/segments/%v/sourceUpdate"
)

type SegmentService struct {
	Config           *types.Config
	FileService      pkg.FileService
	QordobaClient    pkg.QordobaClient
	WorkspaceService pkg.WorkspaceService
}

func (s *SegmentService) AddKey(fileName, version string, keyAddRequest *types.KeyAddRequest) {
	file, _ := s.FileService.FindFile(fileName, version, false)
	if file == nil {
		return
	}
	base := s.Config.GetAPIBase()
	addKeyRequestURL := fmt.Sprintf(keyAddTemplate, base, s.Config.Qordoba.OrganizationID, s.Config.Qordoba.ProjectID, file.FileID)
	log.Debugf("call %v to add key", addKeyRequestURL)
	resp, err := s.QordobaClient.PostToServer(addKeyRequestURL, keyAddRequest)
	if err != nil {
		log.Errorf("error occurred on post key-pair: %v", err)
		return
	}
	handleAddKeyResponse(resp, keyAddRequest, version, fileName)
}

func handleAddKeyResponse(resp *http.Response, keyAddRequest *types.KeyAddRequest, version, fileName string) {
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		if resp.StatusCode == http.StatusUnauthorized {
			log.Errorf("User is not authorised for this request. Check `access_token` in configuration.")
		} else {
			log.Errorf("Problem to add key '%s'. Status: %v\nResponse : %v", keyAddRequest.Key, resp.Status, string(body))
		}
	} else {
		if version == "" {
			log.Infof("Key '%s' was added to file '%s'.", keyAddRequest.Key, fileName)
		} else {
			log.Infof("Key '%s' was added to file '%s' (%v).", keyAddRequest.Key, fileName, version)
		}
	}
}

// UpdateKey function update key
func (s *SegmentService) UpdateKey(fileName, version string, keyAddRequest *types.KeyAddRequest) {
	base := s.Config.GetAPIBase()
	file, personaID := s.FileService.FindFile(fileName, version, false)
	if file == nil {
		return
	}
	segment := s.FindSegment(base, keyAddRequest.Key, personaID, file)
	if segment != nil {
		updateKeyRequestURL := fmt.Sprintf(keyUpdateTemplate, base, s.Config.Qordoba.OrganizationID, s.Config.Qordoba.ProjectID, file.FileID, segment.SegmentID)
		resp, err := s.QordobaClient.PutToServer(updateKeyRequestURL, keyAddRequest)
		handleUpdateKeyResult(resp, err)
	} else {
		if version != "" {
			log.Errorf("Segment %s in %s %s was not found", keyAddRequest.Key, fileName, version)
			return
		}
		log.Errorf("Segment %s in %s was not found", keyAddRequest.Key, fileName)
	}
}

func (s *SegmentService) FindSegment(base, segmentName string, personaID int, file *types.File) *types.Segment {
	workspaceData, err := s.WorkspaceService.LoadWorkspace()
	if err != nil {
		log.Errorf("error occurred on retrieving workspace workspaceData ")
		return nil
	}
	for _, workflow := range workspaceData.Workflow {
		getSegmentRequest := fmt.Sprintf(getSegmentTemplate, base, s.Config.Qordoba.OrganizationID, s.Config.Qordoba.ProjectID, personaID, file.FileID, workflow.ID, segmentName)
		fmt.Printf("segment = %v\n", getSegmentRequest)
		resp, err := s.QordobaClient.GetFromServer(getSegmentRequest)
		if err != nil {
			log.Debugf("error occurred: %v", err)
			continue
		}
		var segmentSearchResponse types.SegmentSearchResponse
		err = json.Unmarshal(resp, &segmentSearchResponse)
		if err != nil {
			log.Errorf("error occurred on server segmentSearchResponse unmarshalling: %v", err)
			continue
		}
		bytes, _ := json.Marshal(segmentSearchResponse)
		log.Debugf("segmentSearch response = %s", string(bytes))
		for _, segment := range segmentSearchResponse.Segments {
			if segment.StringKey == segmentName {
				return &segment
			}
		}
	}
	return nil
}

func handleUpdateKeyResult(resp *http.Response, err error) {
	if err != nil {
		log.Errorf("error occurred on update key attempt: %v", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		if resp.StatusCode == http.StatusUnauthorized {
			log.Errorf("User is not authorised for this request. Check `access_token` in configuration.")
		} else {
			log.Errorf("Segment update status: %v. Response : %v", resp.Status, string(body))
		}
	} else {
		log.Info("Segment was successfully updated")
	}
}

func (s *SegmentService) DeleteKey(fileName, version, segmentKey string) {

}