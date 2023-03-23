package games

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	gameGenreCollection = "game-genres"
)

type GameRepositoryImpl struct {
	mongoDbClient *mongo.Client
	validate      *validator.Validate
}

func NewGameRepositoryImpl(mongoDbClient *mongo.Client) *GameRepositoryImpl {
	return &GameRepositoryImpl{
		mongoDbClient: mongoDbClient,
		validate:      validator.New(),
	}
}

func (g *GameRepositoryImpl) saveGameGenre(ctx context.Context, genre *GameGenre) error {
	_, err := g.mongoDbClient.Database("test").Collection(gameGenreCollection).InsertOne(ctx, genre)

	if err != nil {
		return UnknownError
	}

	return nil
}

func (g *GameRepositoryImpl) updateGameGenre(ctx context.Context, genre *GameGenre) error {

	filter := bson.D{{"slug", genre.Slug}}
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", genre}}

	_, err := g.mongoDbClient.Database("test").Collection(gameGenreCollection).UpdateOne(ctx, filter, update, opts)

	if err != nil {
		return UnknownError
	}

	return nil
}

func (g *GameRepositoryImpl) getGameGenre(ctx context.Context, slug string) (*GameGenre, error) {
	var gameGenre GameGenre

	err := g.mongoDbClient.Database("test").Collection(gameGenreCollection).FindOne(ctx, bson.D{{"slug", slug}}).Decode(&gameGenre)
	if err != nil {
		return nil, ErrNotFound
	}

	return &gameGenre, nil
}

func (g *GameRepositoryImpl) getAllGameGenres(ctx context.Context, pagination *Pagination) (*PaginatedGameGenres, error) {
	var gameGenres []GameGenre

	var response *PaginatedGameGenres

	limit := int64(pagination.Limit)
	skip := int64(pagination.Offset)

	opts := options.Find().SetSort(bson.D{{"dateAdded", -1}, {"updatedAt", -1}}).SetLimit(limit).SetSkip(skip)

	cursor, err := g.mongoDbClient.Database("test").Collection(gameGenreCollection).Find(ctx, bson.D{}, opts)

	if err != nil {
		return nil, UnknownError
	}

	err = cursor.All(ctx, &gameGenres)

	if err != nil {
		return nil, UnknownError
	}

	response = &PaginatedGameGenres{
		Data: gameGenres,
	}

	count, err := g.mongoDbClient.Database("test").Collection(gameGenreCollection).CountDocuments(ctx, bson.D{})

	if err != nil {
		return nil, UnknownError
	}

	response.TotalItems = int(count)
	response.TotalPages = int(count) / pagination.Limit
	response.CurrentPage = pagination.Offset / pagination.Limit
	response.ItemsPerPage = pagination.Limit
	response.HasMore = response.TotalPages > response.CurrentPage

	return response, nil
}

func (g *GameRepositoryImpl) deleteGameGenre(ctx context.Context, slug string) error {
	_, err := g.mongoDbClient.Database("test").Collection(gameGenreCollection).DeleteOne(ctx, bson.D{{"slug", slug}})

	if err != nil {
		return UnknownError
	}

	return nil
}
