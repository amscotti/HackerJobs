package models

// Posting represents a job posting from Hacker News
type Posting struct {
	By          string  `json:"by"`
	Descendants int64   `json:"descendants"`
	ID          int64   `json:"id"`
	Kids        []int64 `json:"kids"`
	Score       int64   `json:"score"`
	Text        string  `json:"text"`
	Time        int64   `json:"time"`
	Title       string  `json:"title"`
	Type        string  `json:"type"`
}
