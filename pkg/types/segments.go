package types

// KeyAddRequest struct for request to add provided key into file
type KeyAddRequest struct {
	Key       string `json:"key"`
	Source    string `json:"source"`
	Reference string `json:"reference"`
}

// ValueKeyUpdateRequest struct for request to add provided key into file
type ValueKeyUpdateRequest struct {
	Segment         string `json:"segment"`
	MoveToFirstStep bool   `json:"moveToFirstStep"`
}

// SegmentSearchResponse struct
type SegmentSearchResponse struct {
	Meta     Meta      `json:"meta"`
	Segments []Segment `json:"segments"`
}

// Segment struct
type Segment struct {
	LastSaved int    `json:"lastSaved"`
	SegmentID int    `json:"segmentId"`
	Reference string `json:"reference"`
	Segment   string `json:"segment"`
	SsText    string `json:"ssText"`
	StringKey string `json:"stringKey"`
}
