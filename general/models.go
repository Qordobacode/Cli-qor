package general

// QordobaConfig structs holds workspace's specific information
type QordobaConfig struct {
	Qordoba Qordoba `yaml:"qordoba"`
	BaseURL string  `yaml:"base_url,omitempty"`
}

type Qordoba struct {
	AccessToken    string `yaml:"access_token"`
	OrganizationID int64  `yaml:"organization_id"`
	ProjectID      int64  `yaml:"project_id"`
}

type PushRequest struct {
	FileName string `json:"filename"`
	Version  string `json:"version"`
	Content  string `json:"content"`
}
