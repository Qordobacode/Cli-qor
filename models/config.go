package models

// QordobaConfig structs holds workspace's specific information
type QordobaConfig struct {
	Qordoba Qordoba `yaml:"qordoba"`
}

type Qordoba struct {
	AccessToken    string `yaml:"access_token"`
	OrganizationID int64  `yaml:"organization_id"`
	ProductID      int64  `yaml:"project_id"`
}
