package games

import "time"

type GameGenre struct {
	Title     string    `json:"title" bson:"title" validate:"required"`
	Slug      string    `json:"slug" bson:"slug" validate:"required"`
	Desc      string    `json:"desc" bson:"desc" validate:"required"`
	CreatedAt time.Time `json:"createdAt" bson:"dateAdded"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	IsDeleted bool      `json:"isDeleted" bson:"isDeleted"`
}
