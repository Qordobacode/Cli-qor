package general

import (
	"github.com/qordobacode/cli-v2/log"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// ApplicationJsonType used in Http header 'Content-Type'
	ApplicationJsonType = "application/json"
)

var (
	// HTTPClient - custom one with a delay set
	HTTPClient = http.Client{
		Timeout: time.Minute * 1,
	}
)

func PostToServer(qordoba *Config, filePath, pushFileURL string, reader io.Reader) {
	request, err := http.NewRequest("POST", pushFileURL, reader)

	if err != nil {
		log.Errorf("error occurred on building file post request: %v", err)
		return
	}
	request.Header.Add("x-auth-token", qordoba.Qordoba.AccessToken)
	request.Header.Add("Content-Type", ApplicationJsonType)
	resp, err := HTTPClient.Do(request)
	if err != nil {
		log.Errorf("error occurred on sending POST request to server")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		log.Errorf("File %s push status: %vresponse : %v", filePath, resp.Status, string(body))
	} else {
		log.Infof("File %s was succesfully pushed to server", filePath)
	}
}
