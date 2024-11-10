package models

type Ad struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	Content     string `json:"content"`
	Impressions int    `json:"impressions"`
}
