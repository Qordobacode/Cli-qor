package rest

import (
	"encoding/json"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server *httptest.Server
)

func buildClient(t *testing.T) *Client {
	// Start a local HTTP server
	server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "access-token", req.Header.Get("x-auth-token"))
		if req.Method == "GET" {
			rw.Write([]byte(`GET RESPONSE`))
		} else if req.Method == "POST" {
			byteResponse, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			var keyAddRequest types.KeyAddRequest
			err = keyAddRequest.UnmarshalJSON(byteResponse)
			assert.Nil(t, err)
			assert.Equal(t, "some-key", keyAddRequest.Key)
			rw.Write([]byte(`POST RESPONSE`))
		} else if req.Method == "PUT" {
			byteResponse, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			var keyAddRequest types.KeyAddRequest
			err = json.Unmarshal(byteResponse, &keyAddRequest)
			assert.Nil(t, err)
			assert.Equal(t, "update-key", keyAddRequest.Key)
			rw.Write([]byte(`PUT RESPONSE`))
		} else if req.Method == "DELETE" {
			// Send response to be tested
			rw.Write([]byte(`DELETE RESPONSE`))
		}
	}))
	return &Client{
		Config: &types.Config{
			Qordoba: types.QordobaConfig{
				AccessToken: "access-token",
			},
		},
		HTTPClient: server.Client(),
	}
}

func TestClientGetFromServer(t *testing.T) {
	client := buildClient(t)
	bytes, err := client.GetFromServer(server.URL)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	assert.Equal(t, "GET RESPONSE", string(bytes))
}

func TestClient_PostToServer(t *testing.T) {
	client := buildClient(t)
	keyAddRequest := types.KeyAddRequest{
		Key: "some-key",
	}
	bytes, err := client.PostToServer(server.URL, keyAddRequest)
	assert.Nil(t, err)
	all, err := ioutil.ReadAll(bytes.Body)

	assert.Equal(t, "POST RESPONSE", string(all))
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
}

func TestClient_PutToServer(t *testing.T) {
	client := buildClient(t)
	keyAddRequest := types.KeyAddRequest{
		Key: "update-key",
	}
	bytes, err := client.PutToServer(server.URL, keyAddRequest)
	assert.Nil(t, err)
	all, err := ioutil.ReadAll(bytes.Body)

	assert.Equal(t, "PUT RESPONSE", string(all))
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
}

func TestClient_DeleteFromServer(t *testing.T) {
	client := buildClient(t)
	bytesResponse, err := client.DeleteFromServer(server.URL)
	assert.Nil(t, err)

	assert.Equal(t, "DELETE RESPONSE", string(bytesResponse))
	assert.Nil(t, err)
}
