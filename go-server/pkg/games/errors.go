package games

import "errors"

var ErrUnauthorized = errors.New("unauthorized")

var ErrBadRequest = errors.New("bad-request")

var ErrNotFound = errors.New("not-found")

var ErrGameGenreAlreadyExists = errors.New("game-genre-already-exists")

var UnknownError = errors.New("internal-server-error")

var ErrGameGenreSlugRequired = errors.New("game-genre-slug-required")

var ErrGameAlreadyExists = errors.New("game-already-exists")
