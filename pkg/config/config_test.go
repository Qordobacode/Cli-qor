package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/qordobacode/cli-v2/pkg/general"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Test struct {
	BaseURL string `mapstructure:"base_url"`
}

func buildConfigService() *ConfigurationService {
	return &ConfigurationService{
		Local: &general.Local{},
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
	service := buildConfigService()
	config, err := service.ReadConfigInPath("config.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "https://app.qordobatest.com/", config.BaseURL)
	assert.Equal(t, "some_token", config.Qordoba.AccessToken)
	assert.Equal(t, int64(3001), config.Qordoba.OrganizationID)
	assert.Equal(t, int64(2879), config.Qordoba.WorkspaceID)
}

func Test_LoadConfig(t *testing.T) {
	service := buildConfigService()
	config, err := service.ReadConfigInPath("config.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "https://app.qordobatest.com/", config.BaseURL)
	assert.Equal(t, "some_token", config.Qordoba.AccessToken)
	assert.Equal(t, int64(3001), config.Qordoba.OrganizationID)
	assert.Equal(t, int64(2879), config.Qordoba.WorkspaceID)
}

func Test_GetConfigPath(t *testing.T) {
	service := buildConfigService()
	configPath, err := service.GetConfigPath()
	assert.Nil(t, err)
	assert.NotNil(t, configPath)
}
