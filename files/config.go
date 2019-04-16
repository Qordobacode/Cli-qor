package files

import (
	"errors"
	"fmt"
	"github.com/qordobacode/cli-v2/log"
	"github.com/qordobacode/cli-v2/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	qordobaHomeTemplate = "%s/.qordoba"
	configPathTemplate  = "%s/.qordoba/config.yaml"
)

func ReadConfigInPath(path string) (*models.QordobaConfig, error) {
	var config models.QordobaConfig
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

// IsConfigFileCorrect validates config file is correct
func IsConfigFileCorrect(config *models.QordobaConfig) bool {
	isConfigCorrect := true
	if config.Qordoba.AccessToken == "" {
		log.Infof("access token is not set\n")
		isConfigCorrect = false
	}
	if config.Qordoba.OrganizationID == 0 {
		log.Infof("organization id is not set\n")
		isConfigCorrect = false
	}
	if config.Qordoba.ProductID == 0 {
		log.Infof("product id is not set\n")
		isConfigCorrect = false
	}
	return isConfigCorrect
}

func IsFilePresent(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// may or may NOT exist. Return false
		return false
	}
}

func PersistAppConfig(config *models.QordobaConfig) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Debugf("error occurred on home dir validation: %v\n", err)
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

func GetConfigPath(home string) (string, error) {
	configPath := fmt.Sprintf(configPathTemplate, home)
	return configPath, nil
}

func IsConfigPresent() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Debugf("error occurred on home dir validation: %v\n", err)
		return false
	}
	configPath, err := GetConfigPath(home)
	if err != nil {
		return false
	}
	return IsFilePresent(configPath)
}
