package config

import (
	"errors"
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

const (
	configName      = ".qordoba.yaml"
	homeConfigName  = "config-v4.yaml"
)

type ConfigurationService struct {
	Local pkg.Local
}

// ReadConfigInPath load config in some folder -> this might be source config OR local config for import
func (c *ConfigurationService) ReadConfigInPath(path string) (*types.Config, error) {
	log.Infof("used config in directory %v", path)
	bytes, err := c.Local.Read(path)
	if err != nil {
		return nil, err
	}
	var config types.Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Errorf("error occurred on config file unmarshaling: %v", err)
		return nil, err
	}
	if !isConfigFileCorrect(&config) {
		return nil, errors.New("config file is incorrect")
	}

	return &config, nil
}

// LoadConfig function loads content of main qordoba configuration
// Read configuration from ~/.qordoba/config-v4.yaml
// Check if current folder contains ./.qordoba.yaml if not search a parent directories for one.
// If you find  set directory with this file as a root to the plugin operations. .qordoba.yaml
// Read content of the  overrides whatever is in  .qordoba.yaml ~/.qordoba/config-v4.yaml
func (c *ConfigurationService) LoadConfig() (*types.Config, error) {
	parentConfig := c.findConfigHierarchically()
	if parentConfig != "" {
		configPath := getConfigPath(parentConfig)
		config, e := c.ReadConfigInPath(configPath)
		if e == nil {
			return config, nil
		}
	}
	return c.readHomeDirectoryConfig()
}

func (c *ConfigurationService) readHomeDirectoryConfig() (*types.Config, error) {
	path, err := c.GetConfigPath()
	if err != nil {
		log.Errorf("error occurred on home dir retrieval: %v", err)
		return nil, err
	}
	return c.ReadConfigInPath(path)
}

func (c *ConfigurationService) findConfigHierarchically() string {
	path, _ := os.Getwd()
	prevPath := path
	for {
		if c.isConfigDir(path) {
			return path
		}
		path = filepath.Clean(path)
		dir, _ := filepath.Split(path)
		if c.isConfigDir(dir) {
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

func (c *ConfigurationService) isConfigDir(path string) bool {
	return c.Local.FileExists(getConfigPath(path))
}

// isConfigFileCorrect validates config file is correct
func isConfigFileCorrect(config *types.Config) bool {
	isConfigCorrect := true
	if config.Qordoba.AccessToken == "" {
		log.Errorf("access_token is not set")
		isConfigCorrect = false
	}
	if config.Qordoba.OrganizationID == 0 {
		log.Errorf("organization_id is not set")
		isConfigCorrect = false
	}
	if config.Qordoba.ProjectID == 0 {
		log.Errorf("product_id is not set")
		isConfigCorrect = false
	}
	return isConfigCorrect
}

// SaveMainConfig function update content of application's config
func (c *ConfigurationService) SaveMainConfig(config *types.Config) {
	marshaledConfig, err := yaml.Marshal(config)
	if err != nil {
		log.Errorf("error occurred on marshalling config file: %v", err)
		return
	}
	c.Local.PutInHome(homeConfigName, marshaledConfig)
}

// GetConfigPath builds path to config according with template
func (c *ConfigurationService) GetConfigPath() (string, error) {
	home, err := c.Local.QordobaHome()
	if err != nil {
		return "", err
	}
	return home + string(os.PathSeparator) + homeConfigName, nil
}
