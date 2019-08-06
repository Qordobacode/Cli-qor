package types

// PushRequest is a request, used to push file
type PushRequest struct {
	FileName string `json:"filename"`
	Version  string `json:"version"`
	Content  string `json:"content"`
}

// WorkspaceResponse is a qordoba general response from obtaining list of workspaces
type WorkspaceResponse struct {
	Meta       Meta            `json:"meta"`
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
	CreatedOn      int64    `json:"createdOn"`
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

// FileSearchResponse - DTO object for.Config's file search response
type FileSearchResponse struct {
	Meta              Meta                `json:"meta"`
	Files             []File              `json:"files"`
	ByPersonaProgress []ByPersonaProgress `json:"byPersonaProgress"`
	TotalCounts       TotalCounts         `json:"totalCounts"`
}

// Paging struct
type Paging struct {
	TotalEnabled int `json:"totalEnabled"`
	TotalResults int `json:"totalResults"`
}

// Meta struct
type Meta struct {
	Paging Paging `json:"paging"`
}

// Tags struct
type Tags struct {
	TagID int    `json:"tagId"`
	Name  string `json:"name"`
}

// Workflow struct contains data about workflow
type Workflow struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Order    int    `json:"order"`
	Complete bool   `json:"complete"`
}

// Counts struct
type Counts struct {
	SegmentCount int `json:"segmentCount"`
	WordCount    int `json:"wordCount"`
}

// ByWorkflowProgress struct
type ByWorkflowProgress struct {
	Workflow Workflow `json:"workflow"`
	Counts   Counts   `json:"counts"`
}

// File struct
type File struct {
	FileID             int                  `json:"fileId"`
	Enabled            bool                 `json:"enabled"`
	Completed          bool                 `json:"completed"`
	Preparing          bool                 `json:"preparing"`
	Filename           string               `json:"filename"`
	Filepath           string               `json:"filepath"`
	Version            string               `json:"version,omitempty"`
	Tags               []Tags               `json:"tags"`
	Update             int64                `json:"update"`
	CreatedAt          int64                `json:"createdAt"`
	Deleted            bool                 `json:"deleted"`
	ByWorkflowProgress []ByWorkflowProgress `json:"byWorkflowProgress"`
	ErrorID            int                  `json:"errorId,omitempty"`
	ErrorMessage       string               `json:"errorMessage,omitempty"`
	Counts             TotalCounts          `json:"counts,omitempty"`
}

// Persona struct contain persona's data
type Persona struct {
	Code string `json:"code"`
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ByPersonaProgress struct used for `status` calculation
type ByPersonaProgress struct {
	Persona            Persona              `json:"persona"`
	ByWorkflowProgress []ByWorkflowProgress `json:"byWorkflowProgress"`
}

// TotalCounts struct
type TotalCounts struct {
	SegmentCount int `json:"segmentCount"`
	WordCount    int `json:"wordCount"`
}

// TagRequest struct
type TagRequest struct {
	TagID int64  `json:"tagId"`
	Name  string `json:"name"`
}

// FileDeleteResponse struct for response for file deletion
type FileDeleteResponse struct {
	Success bool `json:"success"`
}
