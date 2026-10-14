package models

type Game struct {
	ID              int        `json:"_id"`
	Title           string     `json:"title"`
	RequiredAge     string     `json:"required_age"`
	IsFree          bool       `json:"is_free"`
	Description     *string    `json:"description"`
	Image           *string    `json:"image"`
	Website         *string    `json:"website"`
	Price           *float64   `json:"price"`
	Developers      []string   `json:"developers"`
	Publishers      []string   `json:"publishers"`
	Platforms       []string   `json:"platforms"`
	MetacriticScore *int       `json:"metacritic_score"`
	Categories      []Category `json:"categories"`
	Genres          []Genre    `json:"genres"`
	Recommendations *int       `json:"recommendations"`
	ReleaseDate     *string    `json:"release_date"`
	ReleaseYear     *int       `json:"release_year"`
}

type Category struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Genre struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type GamesLinks struct {
	Self string  `json:"self"`
	Next *string `json:"next,omitempty"`
}

type GamesResponse struct {
	Items      []Game     `json:"items"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
	TotalItems int        `json:"totalItems"`
	TotalPages int        `json:"totalPages"`
	Links      GamesLinks `json:"_links"`
}
