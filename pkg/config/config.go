package config

import (
	"errors"
	"github.com/imdario/mergo"
	"github.com/qordobacode/cli-v2/pkg"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

const (
	hiddenConfigName = ".qordoba"
	newConfigName    = "config.yaml"
)

// ConfigurationService is an implementation of pkg.ConfigurationService
type ConfigurationService struct {
	Local pkg.Local
}

// ReadConfigInPath load config in some folder -> this might be source config OR local config for import
func (c *ConfigurationService) ReadConfigInPath(path string) (*types.Config, error) {
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
	return &config, nil
}

// LoadConfig function loads content of main qordoba configuration
// Read configuration from ~/.qordoba/config-v4.yaml
// Check if current folder contains ./.qordoba.yaml if not search a parent directories for one.
// If you find  set directory with this file as a root to the plugin operations. .qordoba.yaml
// Read content of the  overrides whatever is in  .qordoba.yaml ~/.qordoba/config-v4.yaml
func (c *ConfigurationService) LoadConfig() (*types.Config, error) {
	homeDirectoryConfig, homeConfigErr := c.readHomeDirectoryConfig()
	viperConfig, viperErr := c.loadConfigFromViper()
	if homeConfigErr != nil || homeDirectoryConfig == nil {
		if viperErr != nil || viperConfig == nil {
			log.Infof("error on read config file from %v\n%v", viper.ConfigFileUsed(), viperErr)
			os.Exit(1)
		}
		log.Infof("config was taken from %v", viper.ConfigFileUsed())
		validateConfigCorrect(viperConfig)
		return viperConfig, viperErr
	}
	if viperErr != nil || viperConfig == nil {
		log.Infof("config was taken from home directory")
		validateConfigCorrect(homeDirectoryConfig)
		return homeDirectoryConfig, nil
	}
	err := mergo.Merge(viperConfig, *homeDirectoryConfig)
	if err != nil {
		return viperConfig, nil
	}
	validateConfigCorrect(viperConfig)
	log.Infof("merge of configs between '%s' and home directory was used", viper.ConfigFileUsed()) //comment
	return viperConfig, nil
}

func (c *ConfigurationService) loadConfigFromViper() (*types.Config, error) {
	viper.Set("Verbose", true)
	viper.SetConfigName(hiddenConfigName) // name of config file (without extension)
	path, _ := os.Getwd()
	prevPath := path
	for {
		path = filepath.Clean(path)
		viper.AddConfigPath(path)
		dir, _ := filepath.Split(path)
		if prevPath == dir {
			break
		} else {
			prevPath = path
			path = dir
		}
	}
	qordobaHome, err := c.Local.QordobaHome()
	if err == nil {
		viper.AddConfigPath(qordobaHome)
	}
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Debugf("%v", err.Error())
		return nil, err
	}

	var config types.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Errorf("error occurred on unmarshalling properties: %v", err)
	}
	if config.Qordoba.WorkspaceID == 0. && config.Qordoba.ProjectID != 0. {
		log.Infof("field 'project_id' in qordoba configuration was deprecated. Please rename it to 'workspace_id' instead")
		config.Qordoba.WorkspaceID = config.Qordoba.ProjectID
	}
	return &config, err
}

func (c *ConfigurationService) readHomeDirectoryConfig() (*types.Config, error) {
	home, err := c.Local.QordobaHome()
	if err != nil {
		return nil, err
	}
	config, err := c.readConfig(home, "config")
	if config != nil {
		return config, err
	}
	config, err = c.readConfig(home, "config-v4")
	if config != nil {
		log.Infof("Config was taken from 'config-v4.yaml' in home qordoba directory. Please, rename this config file to 'config.yaml'")
		return config, err
	}
	return nil, errors.New("config was not found")
}

func (c *ConfigurationService) readConfig(home, filename string) (*types.Config, error) {
	yamlConfigPath := home + string(os.PathSeparator) + filename + ".yaml"
	if c.Local.FileExists(yamlConfigPath) {
		return c.ReadConfigInPath(yamlConfigPath)
	}
	ymlConfigPath := home + string(os.PathSeparator) + filename + ".yml"
	if c.Local.FileExists(ymlConfigPath) {
		return c.ReadConfigInPath(ymlConfigPath)
	}
	return nil, nil
}

// validateConfigCorrect validates config file is correct
func validateConfigCorrect(config *types.Config) {
	if config == nil {
		log.Errorf("qordoba config was not found")
		os.Exit(1)
	}
	if config.Qordoba.AccessToken == "" {
		log.Errorf("qordoba.access_token is not set")
		os.Exit(1)
	}
	if config.Qordoba.OrganizationID == 0 {
		log.Errorf("qordoba.organization_id is not set")
		os.Exit(1)
	}
	if config.Qordoba.WorkspaceID == 0 {
		log.Errorf("qordoba.workspace_id is not set")
		os.Exit(1)
	}
}

// SaveMainConfig function update content of application's config
func (c *ConfigurationService) SaveMainConfig(config *types.Config) {
	marshaledConfig, err := yaml.Marshal(config)
	if err != nil {
		log.Errorf("error occurred on marshalling config file: %v", err)
		return
	}
	c.Local.PutInHome(newConfigName, marshaledConfig)
}

// GetConfigPath builds path to config according with template
func (c *ConfigurationService) GetConfigPath() (string, error) {
	home, err := c.Local.QordobaHome()
	if err != nil {
		return "", err
	}
	return home + string(os.PathSeparator) + newConfigName, nil
}
