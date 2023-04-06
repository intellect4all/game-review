package reviews

import "errors"

var ErrBadRequest = errors.New("bad-request")

var ErrNotFound = errors.New("not-found")

var UnknownError = errors.New("internal-server-error")

var ErrGameNotFound = errors.New("game-not-found")

var ErrReviewNotFound = errors.New("review-not-found")

var ErrUnauthorized = errors.New("unauthorized")
