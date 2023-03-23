package games

import "errors"

var ErrUnauthorized = errors.New("unauthorized")

var ErrBadRequest = errors.New("bad-request")

var ErrNotFound = errors.New("not-found")

var ErrGameGenreAlreadyExists = errors.New("game-genre-already-exists")

var UnknownError = errors.New("internal-server-error")
