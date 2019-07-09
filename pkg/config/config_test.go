package config

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/mitchellh/mapstructure"
	"github.com/qordobacode/cli-v2/pkg/mock"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

var (
	config = `
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

type Test struct {
	BaseURL string `mapstructure:"base_url"`
}

func buildConfigService(t *testing.T) *ConfigurationService {
	controller := gomock.NewController(t)
	local := mock.NewMockLocal(controller)
	local.EXPECT().QordobaHome().Return("home/qor", nil).Times(2)
	local.EXPECT().FileExists(gomock.Any()).Return(true)
	local.EXPECT().Read(gomock.Any()).Return([]byte(config), nil)
	local.EXPECT().PutInHome(newConfigName, gomock.Any())
	return &ConfigurationService{
		Local: local,
	}
}

func Test_Deccode(t *testing.T) {
	input := map[string]interface{}{
		"base_url": "test",
	}

	var result Test
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", result.BaseURL)
	assert.Equal(t, "test", result.BaseURL)
}

func Test_ReadConfigInPath(t *testing.T) {
	service := buildConfigService(t)
	config, err := service.ReadConfigInPath("config.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "https://app.qordobadev.com/", config.BaseURL)
	assert.Equal(t, "test", config.Qordoba.AccessToken)
	assert.Equal(t, int64(9), config.Qordoba.OrganizationID)
	assert.Equal(t, int64(399), config.Qordoba.WorkspaceID)
}

func Test_LoadConfig(t *testing.T) {
	service := buildConfigService(t)
	config, err := service.ReadConfigInPath("config.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "https://app.qordobadev.com/", config.BaseURL)
	assert.Equal(t, "test", config.Qordoba.AccessToken)
	assert.Equal(t, int64(9), config.Qordoba.OrganizationID)
	assert.Equal(t, int64(399), config.Qordoba.WorkspaceID)
}

func Test_GetConfigPath(t *testing.T) {
	service := buildConfigService(t)
	configPath, err := service.GetConfigPath()
	assert.Nil(t, err)
	assert.NotNil(t, configPath)
}

func TestConfigurationService_LoadConfig(t *testing.T) {
	service := buildConfigService(t)
	config, err := service.LoadConfig()
	assert.Nil(t, err)
	assert.NotNil(t, config)
}

func TestConfigurationService_SaveMainConfig(t *testing.T) {
	service := buildConfigService(t)
	var configObject types.Config
	err := yaml.Unmarshal([]byte(config), &configObject)
	assert.Nil(t, err)
	service.SaveMainConfig(&configObject)
}
