package general

import (
	"fmt"
	"github.com/qordobacode/cli-v2/log"
	"io/ioutil"
	"net/http"
)

var (
	keyAddTemplate     = "%s/v3/organizations/%d/workspaces/%d/files/%d/segments/keyAdd"
	keyUpdateTemplate  = "%s/v3/organizations/%d/workspaces/%d/files/%d/segments/%v/keyUpdate"
	getSegmentTemplate = "%s/v3/organizations/%d/workspaces/%d/personas/%v/files/%d/workflow/%d/segments"
)

func AddKey(config *Config, fileName, version string, keyAddRequest *KeyAddRequest) {
	file := FindFile(config, fileName, version)
	if file == nil {
		return
	}
	base := config.GetAPIBase()
	addKeyRequestURL := fmt.Sprintf(keyAddTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, file.FileID)
	log.Debugf("call %v to add key", addKeyRequestURL)
	resp, err := PostToServer(config, addKeyRequestURL, keyAddRequest)
	if err != nil {
		log.Errorf("error occurred on post key-pair: %v", err)
		return
	}
	handleAddKeyResponse(resp, keyAddRequest, version, fileName)
}

func handleAddKeyResponse(resp *http.Response, keyAddRequest *KeyAddRequest, version, fileName string) {
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
func UpdateKey(config *Config, fileName, version string, keyAddRequest *KeyAddRequest) {
	base := config.GetAPIBase()
	file, personaID := FindFile(config, fileName, version)
	if file == nil {
		return
	}
	segment := FindSegment(config, base, keyAddRequest.Key)
	if segment != nil {
		updateKeyRequestURL := fmt.Sprintf(getSegmentTemplate, base, config.Qordoba.OrganizationID, config.Qordoba.ProjectID, personaID, file.FileID, segment.SegmentID)
		resp, err := PutToServer(config, updateKeyRequestURL, keyAddRequest)
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
		} else {
			log.Errorf("Segment update status: %v. Response : %v", resp.Status, string(body))
		}
	} else {
		log.Info("Segment was succesfully updated")
	}
}

func FindSegment(config *Config, base, segmentName string) *Segment {
	//workspaceData, e := GetWorkspaceData(config)
	//if e != nil {
	//	log.Errorf("error occurred on retrieving workspace workspaceData ")
	//	return nil
	//}
	//for _, workflow := range workspaceData.Workflow {
	//
	//}
	return nil
}
