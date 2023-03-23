package games

type AddGenreRequest struct {
	Title string `json:"title" validate:"required"`
	Desc  string `json:"desc" validate:"required"`
}
