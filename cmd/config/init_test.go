package config

import (
	"github.com/golang/mock/gomock"
	"github.com/qordobacode/cli-v2/pkg/mock"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

var (
	configYAML = `
qordoba:
  access_token: test
  organization_id: 9
  workspace_id: 399
push:
  sources:
    files:
      - C:\data\golang-example\qordoba-example-formats\pot2\core.qordoba-pot2
    folders: []
download:
  targets: []
blacklist:
  sources:
    - ".*.strings"
base_url: "https://app.qordobadev.com/"
`
)

func prepareInit(t *testing.T) {
	controller := gomock.NewController(t)
	service := mock.NewMockConfigurationService(controller)
	var configObject types.Config
	err := yaml.Unmarshal([]byte(configYAML), &configObject)
	assert.Nil(t, err)
	service.EXPECT().ReadConfigInPath("file.txt").Return(&configObject, nil)
	service.EXPECT().SaveMainConfig(gomock.Any())
	configurationService = service
}

func TestConfigFileImported(t *testing.T) {
	initCmd := NewInitCmd()
	prepareInit(t)
	err := RunInitRoot(initCmd, []string{"file.txt"})
	assert.Nil(t, err)
}
