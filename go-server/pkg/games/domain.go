package games

import (
	"context"
	"github.com/go-playground/validator/v10"
	"log"
	"strings"
	"time"
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

func (g *GameService) EditGameGenre(ctx context.Context, genre *GameGenre) error {
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
	// ideally, we should check if the genre is being used by any game before deleting it
	return g.repository.deleteGameGenre(ctx, slug)
}
