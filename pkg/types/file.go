package types

// ScoreResponseBody struct represent score response
type ScoreResponseBody struct {
	SnapshotTime  int64       `json:"snapshotTime"`
	DocumentScore float64     `json:"documentScore"`
	Breakdown     []Breakdown `json:"breakdown"`
}

// Breakdown contains Breakdown information
type Breakdown struct {
	Category   string  `json:"category"`
	IssueCount int     `json:"issueCount"`
	Score      float64 `json:"score"`
	Enabled    bool    `json:"enabled"`
}
