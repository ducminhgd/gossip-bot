package models

import "time"

// News represents a news item
type News struct {
	// Title is the title of the news item
	Title string `json:"title"`
	
	// URL is the URL of the news item
	URL string `json:"url"`
	
	// Description is a short description of the news item
	Description string `json:"description"`
	
	// Source is the source of the news item
	Source string `json:"source"`
	
	// SubSource is the sub-source of the news item (e.g., subreddit name)
	SubSource string `json:"sub_source,omitempty"`
	
	// PublishedAt is the time the news item was published
	PublishedAt time.Time `json:"published_at"`
	
	// Score is the score/rating of the news item
	Score int `json:"score"`
	
	// Comments is the number of comments on the news item
	Comments int `json:"comments"`
}
