package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	// ApplicationJSONType used in Http header 'Content-Type'
	ApplicationJSONType = "application/json"
)

type RestClient struct {
	Config     *types.Config
	HTTPClient http.Client
}

func NewRestClient(qordobaConfig *types.Config) *RestClient {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   true,
	}
	return &RestClient{
		HTTPClient: http.Client{
			Timeout: time.Minute * 1,
			Transport: transport,
		},
		Config: qordobaConfig,
	}
}

// GetFromServer - util function for general request to server. Adds x-auth-token from config, validate response
func (r *RestClient) GetFromServer(getURL string) ([]byte, error) {
	request, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		log.Errorf("error occurred on request build: %v", err)
		return nil, err
	}
	request.Header.Add("x-auth-token", r.Config.Qordoba.AccessToken)
	response, err := r.HTTPClient.Do(request)
	if err != nil {
		log.Errorf("error occurred on workspace get request: %v", err)
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred on body read: %v", err)
	}
	if response.StatusCode/100 != 2 {
		if response.StatusCode == http.StatusUnauthorized {
			log.Errorf("User is not authorised for this request. Check `access_token` in configuration.")
		} else {
			log.Debugf("Error occurred on get %s request. Status: %v, Response : %v", getURL, response.Status, string(bodyBytes))
		}
		return nil, errors.New("unsuccessful request")
	}
	return bodyBytes, err
}

// PostToServer send POST request to server with specified body
func (r *RestClient) PostToServer(postURL string, requestBody interface{}) (*http.Response, error) {
	reader, err := wrapRequest2Reader(requestBody)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", postURL, reader)

	if err != nil {
		return nil, err
	}
	request.Header.Add("x-auth-token", r.Config.Qordoba.AccessToken)
	request.Header.Add("Content-Type", ApplicationJSONType)
	return r.HTTPClient.Do(request)
}

// PutToServer send PUT request to server with specified body
func (r *RestClient) PutToServer(postURL string, requestBody interface{}) (*http.Response, error) {
	reader, err := wrapRequest2Reader(requestBody)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("PUT", postURL, reader)

	if err != nil {
		return nil, err
	}
	request.Header.Add("x-auth-token", r.Config.Qordoba.AccessToken)
	request.Header.Add("Content-Type", ApplicationJSONType)
	return r.HTTPClient.Do(request)
}

func wrapRequest2Reader(requestBody interface{}) (io.Reader, error) {
	marshaledBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Errorf("error occurred on marshalling object: %v", err)
		return nil, err
	}
	reader := bytes.NewReader(marshaledBody)
	return reader, nil
}

// DeleteFromServer - send DELETE request to server
func (r *RestClient) DeleteFromServer(deleteURL string) ([]byte, error) {
	request, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		log.Errorf("error occurred on request build: %v", err)
		return nil, err
	}
	request.Header.Add("x-auth-token", r.Config.Qordoba.AccessToken)
	response, err := r.HTTPClient.Do(request)
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
			log.Errorf("Error occurred on %s request. Status: %d, Response : %v", deleteURL, response.StatusCode, string(bodyBytes))
		}
		return nil, errors.New("unsuccessful request")
	}
	return bodyBytes, err
}
