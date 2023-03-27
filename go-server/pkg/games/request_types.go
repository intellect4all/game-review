package games

import (
	"time"
)

type AddGenreRequest struct {
	Title string `json:"title" validate:"required"`
	Desc  string `json:"desc" validate:"required"`
}

type AddGameRequest struct {
	Title       string               `json:"title" validate:"required"`
	Summary     string               `json:"summary" validate:"required"`
	ReleaseDate time.Time            `json:"releaseDate" validate:"required"`
	Developer   string               `json:"developer" validate:"required"`
	Publisher   string               `json:"publisher" validate:"required"`
	Genres      []*EmbeddedGameGenre `json:"genres" validate:"required"`
}

type UpdateGameRequest struct {
	Title       string               `json:"title" validate:"omitempty"`
	Summary     string               `json:"summary" validate:"omitempty"`
	ReleaseDate time.Time            `json:"releaseDate" validate:"omitempty"`
	Developer   string               `json:"developer" validate:"omitempty"`
	Publisher   string               `json:"publisher" validate:"omitempty"`
	Genres      []*EmbeddedGameGenre `json:"genres" validate:"omitempty"`
}
