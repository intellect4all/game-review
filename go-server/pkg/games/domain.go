package games

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"
	"time"
)

type Service struct {
	repository Repository
	validate   *validator.Validate
}

type PaginatedResponseType interface {
	GameGenre | Game
}

type PaginatedResponse[V PaginatedResponseType] struct {
	Data         []V  `json:"data"`
	CurrentPage  int  `json:"currentPage"`
	TotalPages   int  `json:"totalPages"`
	TotalItems   int  `json:"totalItems"`
	HasMore      bool `json:"hasMore"`
	ItemsPerPage int  `json:"itemsPerPage"`
}

type Pagination struct {
	Limit        int                    `json:"limit"`
	Offset       int                    `json:"offset"`
	QueryFilters map[string]interface{} `json:"filters"`
}

type EmbeddedGameGenre struct {
	Title string `json:"title" bson:"title"`
	Slug  string `json:"slug" bson:"slug"`
}

type RatingStats struct {
	Sum   int `json:"sum" bson:"sum"`
	Count int `json:"count" bson:"count"`
}

type Game struct {
	Title       string               `json:"title" bson:"title"`
	Summary     string               `json:"summary" bson:"summary"`
	Id          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	ReleaseDate time.Time            `json:"releaseDate" bson:"releaseDate"`
	Developer   string               `json:"developer" bson:"developer"`
	Publisher   string               `json:"publisher" bson:"publisher"`
	Genres      []*EmbeddedGameGenre `json:"genres" bson:"genres"`
	Rating      RatingStats          `json:"rating" bson:"rating"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
	IsDeleted   bool                 `json:"isDeleted" bson:"isDeleted"`
	Image       string               `json:"image" bson:"image"`
}

func NewGameService(repository Repository) *Service {
	return &Service{
		validate:   validator.New(),
		repository: repository,
	}
}

type Repository interface {
	saveGameGenre(ctx context.Context, genre *GameGenre) error
	updateGameGenre(ctx context.Context, genre *GameGenre) error
	getGameGenre(ctx context.Context, slug string) (*GameGenre, error)
	getAllGameGenres(ctx context.Context, pagination *Pagination) (*PaginatedResponse[GameGenre], error)
	deleteGameGenre(ctx context.Context, slug string) error
	getGame(ctx context.Context, id string) (*Game, error)
	saveGame(ctx context.Context, game *Game) error
	updateGame(ctx context.Context, game *Game) error
	deleteGame(ctx context.Context, id string) error
	getAllGames(ctx context.Context, pagination *Pagination) (*PaginatedResponse[Game], error)
}

func (g *Service) AddGameGenre(ctx context.Context, genre *GameGenre) error {
	if err := g.validate.Struct(genre); err != nil {
		log.Printf("Validation error: %s", err.Error())
		return err
	}

	_, err := g.repository.getGameGenre(ctx, genre.Slug)

	if err == nil {
		return ErrGameGenreAlreadyExists
	}

	err = g.repository.saveGameGenre(ctx, genre)

	if err != nil {
		return err
	}

	return nil
}

func (g *Service) EditGameGenre(ctx context.Context, genre *GameGenre) error {
	slug := strings.TrimSpace(genre.Slug)

	if genre.Slug == "" {
		return ErrGameGenreSlugRequired
	}

	oldGenre, err := g.repository.getGameGenre(ctx, slug)

	if err != nil {

		return err
	}
	updateGenre(genre, oldGenre)

	err = g.repository.updateGameGenre(ctx, genre)

	if err != nil {
		return err
	}

	return nil
}

func updateGenre(new *GameGenre, old *GameGenre) {
	new.Slug = old.Slug
	new.CreatedAt = old.CreatedAt

	if strings.TrimSpace(new.Title) == "" {
		new.Title = old.Title
	}

	if strings.TrimSpace(new.Desc) == "" {
		new.Desc = old.Desc
	}

	new.UpdatedAt = time.Now()
}

func (g *Service) GetGameGenre(ctx context.Context, slug string) (*GameGenre, error) {

	genre, err := g.repository.getGameGenre(ctx, slug)

	if err != nil {
		return nil, err
	}

	return genre, nil
}

func (g *Service) GetAllGenres(ctx context.Context, pagination *Pagination) (*PaginatedResponse[GameGenre], error) {

	paginatedResponse, err := g.repository.getAllGameGenres(ctx, pagination)

	if err != nil {
		return nil, err
	}

	return paginatedResponse, nil
}

func (g *Service) DeleteGameGenre(ctx context.Context, slug string) error {

	_, err := g.repository.getGameGenre(ctx, slug)

	if err != nil {
		return err
	}
	// ideally, we should isAdminOrModerator if the genre is being used by any game before deleting it
	return g.repository.deleteGameGenre(ctx, slug)
}

func (g *Service) AddGame(ctx context.Context, newGame *Game) error {
	if err := g.validate.Struct(newGame); err != nil {
		return err
	}

	game, err := g.repository.getGame(ctx, newGame.Id.Hex())

	if err == nil && !game.IsDeleted {
		return ErrGameAlreadyExists
	}

	err = g.repository.saveGame(ctx, newGame)

	if err != nil {
		return err
	}

	return nil
}

func (g *Service) UpdateGame(ctx context.Context, game *Game) error {
	if err := g.validate.Struct(game); err != nil {
		return err
	}

	_, err := g.repository.getGame(ctx, game.Id.Hex())

	if err != nil {
		return ErrNotFound
	}

	err = g.repository.updateGame(ctx, game)
	if err != nil {
		return err
	}

	return nil
}

func (g *Service) DeleteGame(ctx context.Context, id string) error {

	err := g.repository.deleteGame(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (g *Service) GetGame(ctx context.Context, id string) (*Game, error) {

	game, err := g.repository.getGame(ctx, id)

	if err != nil {
		return nil, err
	}

	return game, nil
}

func (g *Service) GetAllGames(ctx context.Context, pagination *Pagination) (*PaginatedResponse[Game], error) {
	log.Println("GetAllGames")
	paginatedResponse, err := g.repository.getAllGames(ctx, pagination)

	if err != nil {
		log.Println("GetAllGames error")
		return nil, err
	}

	return paginatedResponse, nil
}
