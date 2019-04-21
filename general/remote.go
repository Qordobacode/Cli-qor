package general

import (
	"errors"
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

// PostToServer send POST request to server with specified body
func PostToServer(qordoba *Config, postURL string, reader io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("POST", postURL, reader)

	if err != nil {
		return nil, err
	}
	request.Header.Add("x-auth-token", qordoba.Qordoba.AccessToken)
	request.Header.Add("Content-Type", ApplicationJsonType)
	return HTTPClient.Do(request)
}

// GetFromServer - util function for general request to server. Adds x-auth-token from config, validate response
func GetFromServer(qordoba *Config, getURL string) ([]byte, error) {
	request, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		log.Errorf("error occurred on request build: %v", err)
		return nil, err
	}
	request.Header.Add("x-auth-token", qordoba.Qordoba.AccessToken)
	response, err := HTTPClient.Do(request)
	if err != nil {
		log.Errorf("error occurred on workspace get request: %v", err)
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("error occurred on body read: %v", err)
		return nil, err
	}
	if response.StatusCode/100 != 2 {
		if response.StatusCode == http.StatusUnauthorized {
			log.Errorf("User is not authorised for this request. Check `access_token` in configuration.")
		} else {
			log.Errorf("Error occurred on %s request. Status: %d, Response : %v", getURL, response.Status, string(bodyBytes))
		}
		return nil, errors.New("unsuccessful request")
	}
	return bodyBytes, err
}
