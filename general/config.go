package general

import (
	"errors"
	"fmt"
	"github.com/qordobacode/cli-v2/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const (
	qordobaHomeTemplate = "%s/.qordoba"
	configPathTemplate  = "%s/.qordoba/config.yaml"
	prodAPIEndpoint     = "https://app.qordoba.com/"
)

// ReadConfigInPath load config in some folder -> this might be source config OR local config for import
func ReadConfigInPath(path string) (*QordobaConfig, error) {
	var config QordobaConfig
	if path == "" {
		log.Infof("Path for config shouldn't be empty\n")
		return nil, errors.New("config path can't be empty")
	}
	// read config from file
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Infof("file not found: %v\n", err)
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Infof("error occurred on config file unmarshaling\n")
		return nil, err
	}
	if !IsConfigFileCorrect(&config) {
		return nil, errors.New("config file is incorrect")
	}

	return &config, nil
}

// LoadConfig function loads content of main quordoba configuration
func LoadConfig() (*QordobaConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Debugf("error occurred on home dir retrieval: %v\n", err)
		return nil, err
	}
	path, err := GetConfigPath(home)
	if err != nil {
		log.Debugf("error occurred on building path to config: %v\n", err)
		return nil, err
	}
	return ReadConfigInPath(path)
}

// IsConfigFileCorrect validates config file is correct
func IsConfigFileCorrect(config *QordobaConfig) bool {
	isConfigCorrect := true
	if config.Qordoba.AccessToken == "" {
		log.Infof("access_token is not set\n")
		isConfigCorrect = false
	}
	if config.Qordoba.OrganizationID == 0 {
		log.Infof("organization_id is not set\n")
		isConfigCorrect = false
	}
	if config.Qordoba.ProjectID == 0 {
		log.Infof("product_id is not set\n")
		isConfigCorrect = false
	}
	return isConfigCorrect
}

// SaveMainConfig function update content of application's config
func SaveMainConfig(config *QordobaConfig) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Debugf("error occurred on home dir retrieval: %v\n", err)
		return
	}
	path, e := GetConfigPath(home)
	if e != nil {
		return
	}
	marshaledConfig, err := yaml.Marshal(config)
	if err != nil {
		log.Infof("error occurred on marshalling config file: %v\n", err)
		return
	}
	qordobaHome := fmt.Sprintf(qordobaHomeTemplate, home)
	err = os.MkdirAll(qordobaHome, os.ModePerm)
	if err != nil {
		log.Infof("error occurred on creating qordoba's folder: %v\n", err)
	}
	err = ioutil.WriteFile(path, marshaledConfig, 0644)
	if err != nil {
		log.Infof("error occurred on writing config: %v\n", err)
	}
}

// GetConfigPath builds path to config according with template
func GetConfigPath(home string) (string, error) {
	configPath := fmt.Sprintf(configPathTemplate, home)
	return configPath, nil
}

// GetAPIBase get value of API endpoint from config OR prod as a default
func (config *QordobaConfig) GetAPIBase() string {
	base := prodAPIEndpoint
	if config.BaseURL != "" {
		base = config.BaseURL
	}
	base = strings.TrimSuffix(base, "/")
	return base
}
