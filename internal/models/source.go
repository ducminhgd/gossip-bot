package models

// Source represents a news source
type Source struct {
	// Name is the name of the source
	Name string `json:"name"`
	
	// Type is the type of the source (e.g., "hackernews", "reddit")
	Type string `json:"type"`
	
	// URL is the base URL of the source
	URL string `json:"url"`
	
	// Limit is the maximum number of news items to fetch
	Limit int `json:"limit"`
	
	// SubSource is the sub-source for sources like Reddit (e.g., "r/golang")
	SubSource string `json:"sub_source,omitempty"`
}
