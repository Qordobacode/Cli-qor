package types

// KeyAddRequest struct for request to add provided key into file
type KeyAddRequest struct {
	Key       string `json:"key"`
	Source    string `json:"source"`
	Reference string `json:"reference"`
}

type SegmentSearchResponse struct {
	Meta     Meta      `json:"meta"`
	Segments []Segment `json:"segments"`
}

type Segment struct {
	LastSaved  int    `json:"lastSaved"`
	Order      int    `json:"order"`
	PluralRule string `json:"pluralRule"`
	Plurals    string `json:"plurals"`
	Reference  string `json:"reference"`
	Segment    string `json:"segment"`
	SegmentID  int    `json:"segmentId"`
	SsMatch    int    `json:"ssMatch"`
	SsText     string `json:"ssText"`
	StringKey  string `json:"stringKey"`
	Target     string `json:"target"`
	TargetID   int    `json:"targetId"`
}
