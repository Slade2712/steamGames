package models

type Game struct {
	ID              int      `json:"_id"`
	Title           string   `json:"title"`
	RequiredAge     int      `json:"required_age"`
	IsFree          bool     `json:"is_free"`
	Description     string   `json:"description"`
	Image           string   `json:"image"`
	Website         string   `json:"website"`
	Developers      []string `json:"developers"`
	Publishers      []string `json:"publishers"`
	Platforms       []string `json:"platforms"`
	Price           *float64 `json:"price"`
	MetacriticScore *int     `json:"metacritic_score"`
	Categories      []Tag    `json:"categories"`
	Genres          []Tag    `json:"genres"`
	Recommendations int      `json:"recommendations"`
	ReleaseDate     string   `json:"release_date"`
	ReleaseYear     int      `json:"release_year"`
}

type Tag struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}
