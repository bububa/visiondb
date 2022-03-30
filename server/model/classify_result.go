package model

// ClassifyResult represents classify results with matched score
type ClassifyResult struct {
	ID     int       `json:"id"`
	Name   string    `json:"name"`
	Score  float64   `json:"score,omitempty"`
	Scores []float64 `json:"scores,omitempty"`
}
