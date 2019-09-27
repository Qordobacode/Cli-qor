package segments

import (
	"fmt"
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	keyAddTemplate     = "%s/v3/organizations/%d/workspaces/%d/files/%d/segments/keyAdd"
	getSegmentTemplate = "%s/v3/organizations/%d/workspaces/%d/personas/%v/files/%d/workflow/%d/segments?search=%s"
	keyUpdateTemplate  = "%s/v3/organizations/%d/workspaces/%d/files/%d/segments/%v/sourceUpdate"
	keyDeleteTemplate  = "%s/v3/organizations/%d/workspaces/%d/files/%d/segments/%v/keyDelete"
)

// SegmentService struct is an implementation of pkg.SegmentService
type SegmentService struct {
	Config           *types.Config
	FileService      pkg.FileService
	QordobaClient    pkg.QordobaClient
	WorkspaceService pkg.WorkspaceService
}

// AddKey function add new key into file
func (s *SegmentService) AddKey(fileName, version string, keyAddRequest *types.KeyAddRequest) {
	keyAddRequest.Key = s.handleSegmentKey(keyAddRequest.Key)
	file, _ := s.FileService.FindFile(fileName, version, false)
	if file == nil {
		return
	}
	base := s.Config.GetAPIBase()
	addKeyRequestURL := fmt.Sprintf(keyAddTemplate, base, s.Config.Qordoba.OrganizationID, s.Config.Qordoba.WorkspaceID, file.FileID)
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
			os.Exit(1)
		} else if resp.StatusCode == http.StatusNotAcceptable {
			log.Errorf("Problem to add key '%s'. Key already exist", keyAddRequest.Key)
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
	keyAddRequest.Key = s.handleSegmentKey(keyAddRequest.Key)
	segment, file := s.FindSegment(fileName, version, keyAddRequest.Key)
	if segment != nil {
		base := s.Config.GetAPIBase()
		updateKeyRequestURL := fmt.Sprintf(keyUpdateTemplate, base, s.Config.Qordoba.OrganizationID, s.Config.Qordoba.WorkspaceID, file.FileID, segment.SegmentID)
		valueUpdateRequest := &types.ValueKeyUpdateRequest{
			Segment:         keyAddRequest.Source,
			MoveToFirstStep: false,
		}
		resp, err := s.QordobaClient.PutToServer(updateKeyRequestURL, valueUpdateRequest)
		handleUpdateKeyResult(resp, err)
	}
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
			os.Exit(1)
		} else {
			log.Errorf("Error on update key: %v. Response: %v", resp.Status, string(body))
		}
	} else {
		log.Info("Segment was successfully updated")
	}
}

// DeleteKey deletes segment from file by key
func (s *SegmentService) DeleteKey(fileName, version, segmentKey string) {
	segmentKey = s.handleSegmentKey(segmentKey)
	segment, file := s.FindSegment(fileName, version, segmentKey)
	if segment != nil {
		base := s.Config.GetAPIBase()
		updateKeyRequestURL := fmt.Sprintf(keyDeleteTemplate, base, s.Config.Qordoba.OrganizationID, s.Config.Qordoba.WorkspaceID, file.FileID, segment.SegmentID)
		_, err := s.QordobaClient.DeleteFromServer(updateKeyRequestURL)
		if err == nil {
			if version != "" {
				log.Infof("Segment %v was successfully deleted from %s - %s", segmentKey, fileName, version)
			} else {
				log.Infof("Segment %v was successfully deleted from %s", segmentKey, fileName)
			}
		}
	}
}

func (s *SegmentService) handleSegmentKey(segmentKey string) string {
	splittedKeys := strings.Split(segmentKey, "/")
	if len(splittedKeys) == 1 {
		log.Info(`Please add "/" to the start of the key`)
		os.Exit(1)
	}
	key := splittedKeys[len(splittedKeys)-1]
	return "/" + key
}

// FindSegment returns segment and file where it is placed by file name/version and segment key
func (s *SegmentService) FindSegment(fileName, fileVersion, key string) (*types.Segment, *types.File) {
	base := s.Config.GetAPIBase()
	file, personaID := s.FileService.FindFile(fileName, fileVersion, false)
	if file == nil {
		return nil, nil
	}
	segment := s.findFileSegment(base, key, personaID, file)
	if segment == nil {
		if fileVersion != "" {
			log.Errorf("Segment %s in %s - %s was not found", key, fileName, fileVersion)
		} else {
			log.Errorf("Segment %s in %s was not found", key, fileName)
		}
	}
	return segment, file
}

func (s *SegmentService) findFileSegment(base, segmentName string, personaID int, file *types.File) *types.Segment {
	workspaceData, err := s.WorkspaceService.LoadWorkspace()
	if err != nil {
		return nil
	}
	for _, workflow := range workspaceData.Workflow {
		getSegmentRequest := fmt.Sprintf(getSegmentTemplate, base, s.Config.Qordoba.OrganizationID, s.Config.Qordoba.WorkspaceID, personaID, file.FileID, workflow.ID, segmentName)
		resp, err := s.QordobaClient.GetFromServer(getSegmentRequest)
		if err != nil {
			log.Debugf("error occurred: %v", err)
			continue
		}
		var segmentSearchResponse types.SegmentSearchResponse
		err = segmentSearchResponse.UnmarshalJSON(resp)
		if err != nil {
			log.Errorf("error occurred on server segmentSearchResponse unmarshalling: %v", err)
			continue
		}
		for _, segment := range segmentSearchResponse.Segments {
			if segment.StringKey == segmentName {
				return &segment
			}
		}
	}
	return nil
}
