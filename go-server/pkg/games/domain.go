package games

import (
	"context"
	"github.com/go-playground/validator/v10"
	"log"
)

type GameService struct {
	repository GameRepository
	validate   *validator.Validate
}

type PaginatedGameGenres struct {
	Data         []GameGenre `json:"data"`
	CurrentPage  int         `json:"currentPage"`
	TotalPages   int         `json:"totalPages"`
	TotalItems   int         `json:"totalItems"`
	HasMore      bool        `json:"hasMore"`
	ItemsPerPage int         `json:"itemsPerPage"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func NewGameService(repository GameRepository) *GameService {
	return &GameService{
		validate:   validator.New(),
		repository: repository,
	}
}

type GameRepository interface {
	saveGameGenre(ctx context.Context, genre *GameGenre) error
	updateGameGenre(ctx context.Context, genre *GameGenre) error
	getGameGenre(ctx context.Context, slug string) (*GameGenre, error)
	getAllGameGenres(ctx context.Context, pagination *Pagination) (*PaginatedGameGenres, error)
	deleteGameGenre(ctx context.Context, slug string) error
}

func (g *GameService) AddGameGenre(ctx context.Context, genre *GameGenre) error {
	log.Println("Authenticated")

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

	log.Println("Game genre added successfully")

	return nil
}

func (g *GameService) EditGameGenre(ctx context.Context, genre *GameGenre) error {

	if err := g.validate.Struct(genre); err != nil {
		return err
	}
	_, err := g.repository.getGameGenre(ctx, genre.Slug)

	if err != nil {
		return err
	}

	err = g.repository.updateGameGenre(ctx, genre)

	if err != nil {
		return err
	}

	return nil
}

func (g *GameService) GetGameGenre(ctx context.Context, slug string) (*GameGenre, error) {

	genre, err := g.repository.getGameGenre(ctx, slug)

	if err != nil {
		return nil, err
	}

	return genre, nil
}

func (g *GameService) GetAllGenres(ctx context.Context, pagination *Pagination) (*PaginatedGameGenres, error) {

	paginatedResponse, err := g.repository.getAllGameGenres(ctx, pagination)

	if err != nil {
		return nil, err
	}

	return paginatedResponse, nil
}

func (g *GameService) DeleteGameGenre(ctx context.Context, slug string) error {

	_, err := g.repository.getGameGenre(ctx, slug)

	if err != nil {
		return err
	}

	return g.repository.deleteGameGenre(ctx, slug)
}
