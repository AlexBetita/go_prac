package models

type BotResponse struct {
	Type     string `json:"type"`     // e.g., "post", "interaction", "related_posts"
	Response any    `json:"response"` // holds any type: *Post, []*Post, *Interaction, etc.
}