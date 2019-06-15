package types

type ScoreResponseBody struct {
	SnapshotTime  int64       `json:"snapshotTime"`
	DocumentScore float64     `json:"documentScore"`
	Breakdown     []Breakdown `json:"breakdown"`
}
type Breakdown struct {
	Category   string  `json:"category"`
	IssueCount int     `json:"issueCount"`
	Score      float64 `json:"score"`
	Enabled    bool    `json:"enabled"`
}
