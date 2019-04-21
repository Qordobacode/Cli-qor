package general

import (
	"io"
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

func GetFromServer(qordoba *Config, pushFileURL string) (*http.Response, error) {
	request, err := http.NewRequest("GET", pushFileURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("x-auth-token", qordoba.Qordoba.AccessToken)
	return HTTPClient.Do(request)
}
