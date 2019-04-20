package general

import (
	"errors"
	"github.com/qordobacode/cli-v2/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	prodAPIEndpoint = "https://app.qordoba.com/"
	configName      = ".qordoba.yaml"
	homeConfigName  = "config-v4.yaml"
)

// ReadConfigInPath load config in some folder -> this might be source config OR local config for import
func ReadConfigInPath(path string) (*Config, error) {
	log.Debugf("used config in directory %v", path)
	var config Config
	if path == "" {
		log.Infof("Path for config shouldn't be empty")
		return nil, errors.New("config path can't be empty")
	}
	// read config from file
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Infof("file not found: %v", err)
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Infof("error occurred on config file unmarshaling")
		return nil, err
	}
	if !IsConfigFileCorrect(&config) {
		return nil, errors.New("config file is incorrect")
	}

	return &config, nil
}

// LoadConfig function loads content of main quordoba configuration
// Read configuration from ~/.qordoba/config-v4.yaml
// Check if current folder contains ./.qordoba.yaml if not search a parent directories for one.
// If you find  set directory with this file as a root to the plugin operations. .qordoba.yaml
// Read content of the  overrides whatever is in  .qordoba.yaml ~/.qordoba/config-v4.yaml
func LoadConfig() (*Config, error) {
	parentConfig := findConfigHierarchically()
	if parentConfig != "" {
		configPath := getConfigPath(parentConfig)
		config, e := ReadConfigInPath(configPath)
		if e == nil {
			return config, nil
		}
	}
	return readHomeDirectoryConfig()
}

func readHomeDirectoryConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Debugf("error occurred on home dir retrieval: %v", err)
		return nil, err
	}
	path := GetConfigPath(home)
	return ReadConfigInPath(path)
}

func findConfigHierarchically() string {
	path, _ := os.Getwd()
	prevPath := path
	for {
		if isConfigDir(path) {
			return path
		}
		path = filepath.Clean(path)
		dir, _ := filepath.Split(path)
		if isConfigDir(dir) {
			return dir
		} else if prevPath == dir {
			return ""
		} else {
			prevPath = path
			path = dir
		}
	}
}

func getConfigPath(path string) string {
	return filepath.Join(path, configName)
}

func isConfigDir(path string) bool {
	return FileExists(getConfigPath(path))
}

// FileExists checks for file existence
func FileExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

// IsConfigFileCorrect validates config file is correct
func IsConfigFileCorrect(config *Config) bool {
	isConfigCorrect := true
	if config.Qordoba.AccessToken == "" {
		log.Infof("access_token is not set")
		isConfigCorrect = false
	}
	if config.Qordoba.OrganizationID == 0 {
		log.Infof("organization_id is not set")
		isConfigCorrect = false
	}
	if config.Qordoba.ProjectID == 0 {
		log.Infof("product_id is not set")
		isConfigCorrect = false
	}
	return isConfigCorrect
}

// SaveMainConfig function update content of application's config
func SaveMainConfig(config *Config) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Debugf("error occurred on home dir retrieval: %v", err)
		return
	}
	path := GetConfigPath(home)
	marshaledConfig, err := yaml.Marshal(config)
	if err != nil {
		log.Infof("error occurred on marshalling config file: %v", err)
		return
	}
	qordobaHome := getQordobaHomeDir(home)
	err = os.MkdirAll(qordobaHome, os.ModePerm)
	if err != nil {
		log.Infof("error occurred on creating qordoba's folder: %v", err)
	}
	err = ioutil.WriteFile(path, marshaledConfig, 0644)
	if err != nil {
		log.Infof("error occurred on writing config: %v", err)
	}
}

// GetConfigPath builds path to config according with template
func GetConfigPath(home string) string {
	return getQordobaHomeDir(home) + string(os.PathSeparator) + homeConfigName
}

func getQordobaHomeDir(home string) string {
	return home + string(os.PathSeparator) + ".qordoba"
}

// GetAPIBase get value of API endpoint from config OR prod as a default
func (config *Config) GetAPIBase() string {
	base := prodAPIEndpoint
	if config.BaseURL != "" {
		base = config.BaseURL
	}
	base = strings.TrimSuffix(base, "/")
	return base
}
