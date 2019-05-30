package types

import "strings"

const (
	prodAPIEndpoint = "https://app.qordoba.com/"
)

// Config structs holds workspace's specific information
type Config struct {
	Qordoba   QordobaConfig   `yaml:"qordoba"`
	Push      PushConfig      `yaml:"push"`
	Download  DownloadConfig  `yaml:"download"`
	Blacklist BlacklistConfig `yaml:"blacklist"`
	BaseURL   string          `yaml:"base_url,omitempty"`
}

// QordobaConfig is a part of configuration with qordoba-related information
type QordobaConfig struct {
	AccessToken    string            `yaml:"access_token"`
	OrganizationID int64             `yaml:"organization_id"`
	ProjectID      int64             `yaml:"project_id"`
	AudienceMap    map[string]string `yaml:"audiences_map"`
}

// PushConfig is push-related part of config
type PushConfig struct {
	Sources SourceConfig `yaml:"sources"`
}

// SourceConfig contains details about source configuration for push config
type SourceConfig struct {
	Files   []string `yaml:"files"`
	Folders []string `yaml:"folders"`
}

// DownloadConfig is download-related part of config
type DownloadConfig struct {
	Targets []string `yaml:"targets"`
}

// BlacklistConfig is blacklist-related part of config
type BlacklistConfig struct {
	Sources []string `yaml:"sources"`
}

// GetAPIBase get value of API endpoint from config OR prod as a default
func (c *Config) GetAPIBase() string {
	base := prodAPIEndpoint
	if c.BaseURL != "" {
		base = c.BaseURL
	}
	base = strings.TrimSuffix(base, "/")
	return base
}

// Audiences function retrieves all languages from audience map
func (c *Config) Audiences() map[string]bool {
	results := make(map[string]bool)
	for _, lang := range c.Qordoba.AudienceMap {
		results[lang] = true
	}
	return results
}
