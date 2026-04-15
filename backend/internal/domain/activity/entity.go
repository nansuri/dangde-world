package activity

import "dangde-world/backend/internal/domain/category"

type Activity struct {
	ID          uint              `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Language    string            `json:"language"`
	Difficulty  string            `json:"difficulty"`
	AgeGroup    string            `json:"ageGroup"`
	CategoryID  uint              `json:"categoryId"`
	Category    category.Category `json:"category"`
	Prompt      string            `json:"prompt"`
	Icon        string            `json:"icon"`
	HTMLCode    string            `json:"htmlCode"`
	CSSCode     string            `json:"cssCode"`
	JSCode      string            `json:"jsCode"`
}
