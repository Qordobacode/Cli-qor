package general

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

// PushRequest is a request, used to push file
type PushRequest struct {
	FileName string `json:"filename"`
	Version  string `json:"version"`
	Content  string `json:"content"`
}

// WorkspaceResponse is a qordoba general response from obtaining list of workspaces
type WorkspaceResponse struct {
	Meta struct {
		Paging struct {
			TotalResults int `json:"totalResults"`
		} `json:"paging"`
	} `json:"meta"`
	Workspaces []WorkspaceData `json:"workspaces"`
}

// WorkspaceData contains workflow and workspace data
type WorkspaceData struct {
	Workflow  []Workflow `json:"workflow"`
	Workspace Workspace  `json:"workspace"`
}

// Workspace is qordoba object with workspace's parameters
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

// Person - qordoba's response with person's information
type Person struct {
	Code      string `json:"code"`
	Direction string `json:"direction"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
}

// FileSearchResponse - DTO object for qordoba's file search response
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
	FileID             int                  `json:"fileId"`
	Enabled            bool                 `json:"enabled"`
	Completed          bool                 `json:"completed"`
	Preparing          bool                 `json:"preparing"`
	Filename           string               `json:"filename"`
	Filepath           string               `json:"filepath"`
	Version            string               `json:"version"`
	Tags               []Tags               `json:"tags"`
	Update             int64                `json:"update"`
	CreatedAt          int64                `json:"createdAt"`
	ErrorID            int                  `json:"errorId"`
	ErrorMessage       string               `json:"errorMessage"`
	Deleted            bool                 `json:"deleted"`
	ByWorkflowProgress []ByWorkflowProgress `json:"byWorkflowProgress"`
}
type Tags struct {
	TagID int    `json:"tagId"`
	Name  string `json:"name"`
}
type Workflow struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Order    int    `json:"order"`
	Complete bool   `json:"complete"`
}

type Counts struct {
	SegmentCount int `json:"segmentCount"`
	WordCount    int `json:"wordCount"`
}
type ByWorkflowProgress struct {
	Workflow Workflow `json:"workflow"`
	Counts   Counts   `json:"counts"`
}

type TagRequest struct {
	TagID int64  `json:"tagId"`
	Name  string `json:"name"`
}

// FileDeleteResponse struct for response for file deletion
type FileDeleteResponse struct {
	Success bool `json:"success"`
}

// KeyAddRequest struct for request to add provided key into file
type KeyAddRequest struct {
	Key       string
	Source    string
	Reference string
}
