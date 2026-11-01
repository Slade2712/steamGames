package games

import "errors"

var (
	ErrInvalidGameID      = errors.New("invalid game id")
	ErrGameNotFound       = errors.New("game not found")
	ErrInvalidSkip        = errors.New("invalid skip")
	ErrInvalidLimit       = errors.New("invalid limit")
	ErrInvalidReleaseYear = errors.New("invalid release_year")
	ErrInvalidOrdering    = errors.New("invalid ordering")
	ErrInvalidGenre       = errors.New("invalid genre")
	ErrInvalidGameData = errors.New("invalid game data")
)
