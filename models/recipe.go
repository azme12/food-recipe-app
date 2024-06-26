package models

type Recipe struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	Steps       []string `json:"steps"`
	PrepTime    int      `json:"time"`
	CategoryID  int64    `json:"category_id"`
	CreatorID   int64    `json:"creator_id"`
	Images      []string `json:"images"`
}