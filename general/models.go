package general

// Config structs holds workspace's specific information
type Config struct {
	Qordoba   QordobaConfig   `yaml:"qordoba"`
	Push      PushConfig      `yaml:"push"`
	Download  DownloadConfig  `yaml:"download"`
	Blacklist BlacklistConfig `yaml:"blacklist"`
	BaseURL   string          `yaml:"base_url,omitempty"`
}

type QordobaConfig struct {
	AccessToken    string            `yaml:"access_token"`
	OrganizationID int64             `yaml:"organization_id"`
	ProjectID      int64             `yaml:"project_id"`
	AudienceMap    map[string]string `yaml:"audiences_map"`
}

type PushConfig struct {
	Sources SourceConfig `yaml:"sources"`
}

type SourceConfig struct {
	Files   []string `yaml:"files"`
	Folders []string `yaml:"folders"`
}

type DownloadConfig struct {
	Targets []string `yaml:"targets"`
}

type BlacklistConfig struct {
	Sources []string `yaml:"sources"`
}

type PushRequest struct {
	FileName string `json:"filename"`
	Version  string `json:"version"`
	Content  string `json:"content"`
}

type WorkspaceResponse struct {
	Meta struct {
		Paging struct {
			TotalResults int `json:"totalResults"`
		} `json:"paging"`
	} `json:"meta"`
	Workspaces []WorkspaceData `json:"workspaces"`
}

type WorkspaceData struct {
	Workflow []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Order int    `json:"order"`
	} `json:"workflow"`
	Workspace Workspace `json:"workspace"`
}

type Workspace struct {
	ContentTypeCodes []interface{} `json:"contentTypeCodes"`
	CreatedBy        struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	} `json:"createdBy"`
	CreatedOn      int      `json:"createdOn"`
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	OrganizationID int      `json:"organizationId"`
	Segmentation   string   `json:"segmentation"`
	SourcePersona  Person   `json:"sourcePersona"`
	TargetPersonas []Person `json:"targetPersonas"`
	Timezone       string   `json:"timezone"`
	TmMatchMode    string   `json:"tmMatchMode"`
}

type Person struct {
	Code      string `json:"code"`
	Direction string `json:"direction"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
}

type FileSearchResponse struct {
	Files []File `json:"files"`
	Meta  struct {
		Paging struct {
			TotalEnabled int `json:"totalEnabled"`
			TotalResults int `json:"totalResults"`
		} `json:"paging"`
	} `json:"meta"`
}

type File struct {
	Completed    bool   `json:"completed"`
	CreatedAt    int64  `json:"createdAt"`
	Deleted      bool   `json:"deleted"`
	Enabled      bool   `json:"enabled"`
	ErrorID      int    `json:"errorId"`
	ErrorMessage string `json:"errorMessage"`
	FileID       int    `json:"fileId"`
	Filename     string `json:"filename"`
	Filepath     string `json:"filepath"`
	Preparing    bool   `json:"preparing"`
	Update       int64  `json:"update"`
	URL          string `json:"url"`
	Version      string `json:"version"`
}

type FileDeleteResponse struct {
	Success bool `json:"success"`
}
