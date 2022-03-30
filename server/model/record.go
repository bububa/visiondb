package model

// Record represents database record
type Record struct {
	// ID Record ID
	ID int `json:"id"`
	// Name Record Name
	Name string `json:"name,omitempty"`
	// ItemsCount number of items
	ItemsCount int `json:"items_count,omitempty"`
}
