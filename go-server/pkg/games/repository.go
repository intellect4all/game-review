package games

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	gameGenreCollection = "genres"
	gamesCollection     = "games"
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
	filter := bson.D{{"slug", slug}, {"isDeleted", false}}

	err := g.mongoDbClient.Database("test").Collection(gameGenreCollection).FindOne(ctx, filter).Decode(&gameGenre)
	if err != nil {
		return nil, ErrNotFound
	}

	return &gameGenre, nil
}

func (g *GameRepositoryImpl) getAllGameGenres(ctx context.Context, pagination *Pagination) (*PaginatedResponse[GameGenre], error) {
	var gameGenres []GameGenre

	response := &PaginatedResponse[GameGenre]{}

	limit := int64(pagination.Limit)
	skip := int64(pagination.Offset)

	filter := bson.D{{"isDeleted", false}}
	opts := options.Find().SetSort(bson.D{{"dateAdded", -1}}).SetLimit(limit).SetSkip(skip)

	cursor, err := g.mongoDbClient.Database("test").Collection(gameGenreCollection).Find(ctx, filter, opts)

	if err != nil {
		return nil, UnknownError
	}

	err = cursor.All(ctx, &gameGenres)

	if err != nil {
		return nil, UnknownError
	}

	if len(gameGenres) == 0 {
		return &PaginatedResponse[GameGenre]{
			TotalItems:   0,
			TotalPages:   0,
			CurrentPage:  0,
			ItemsPerPage: 0,
			HasMore:      false,
			Data:         []GameGenre{},
		}, nil
	}

	response.Data = gameGenres

	count, err := g.mongoDbClient.Database("test").Collection(gameGenreCollection).CountDocuments(ctx, bson.D{}, options.Count())

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
	filter := bson.D{{"slug", slug}}

	update := bson.D{{"$set", bson.D{{"isDeleted", true}}}}

	opts := options.Update().SetUpsert(true)

	_, err := g.mongoDbClient.Database("test").Collection(gameGenreCollection).UpdateOne(ctx, filter, update, opts)

	if err != nil {
		return UnknownError
	}

	return nil
}

func (g *GameRepositoryImpl) getGame(ctx context.Context, id string) (*Game, error) {
	var game Game

	err := g.mongoDbClient.Database("test").Collection(gamesCollection).FindOne(ctx, bson.D{{"_id", id}}).Decode(&game)
	if err != nil {
		return nil, ErrNotFound
	}

	return &game, nil
}

func (g *GameRepositoryImpl) getAllGames(ctx context.Context, pagination *Pagination) (*PaginatedResponse[Game], error) {

	var games []Game

	response := &PaginatedResponse[Game]{}

	limit := int64(pagination.Limit)
	skip := int64(pagination.Offset)

	opts := options.Find().SetSort(bson.D{{"createdAt", -1}, {"updatedAt", -1}}).SetLimit(limit).SetSkip(skip)

	filter := bson.D{{"isDeleted", false}}

	if pagination.QueryFilters != nil {

		for key, value := range pagination.QueryFilters {
			log.Println("Key: ", key, " Value: ", value)
			filter = append(filter, bson.E{Key: key, Value: value})
		}
	}

	cursor, err := g.mongoDbClient.Database("test").Collection(gamesCollection).Find(ctx, filter, opts)

	if err != nil {

		return nil, UnknownError
	}

	err = cursor.All(ctx, &games)

	if err != nil {

		return nil, UnknownError
	}

	if len(games) == 0 {
		return &PaginatedResponse[Game]{
			TotalItems:   0,
			TotalPages:   0,
			CurrentPage:  0,
			ItemsPerPage: 0,
			HasMore:      false,
			Data:         []Game{},
		}, nil
	}

	response.Data = games

	count, err := g.mongoDbClient.Database("test").Collection(gamesCollection).CountDocuments(ctx, filter, options.Count())

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

func (g *GameRepositoryImpl) saveGame(ctx context.Context, game *Game) error {
	_, err := g.mongoDbClient.Database("test").Collection(gamesCollection).InsertOne(ctx, game)

	if err != nil {
		return UnknownError
	}

	return nil
}

func (g *GameRepositoryImpl) updateGame(ctx context.Context, game *Game) error {

	filter := bson.D{{"_id", game.Id}}
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", game}}

	_, err := g.mongoDbClient.Database("test").Collection(gamesCollection).UpdateOne(ctx, filter, update, opts)

	if err != nil {
		return UnknownError
	}

	return nil
}

func (g *GameRepositoryImpl) deleteGame(ctx context.Context, id string) error {
	filter := bson.D{{"_id", id}}
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", bson.D{{"IsDeleted", true}}}}

	_, err := g.mongoDbClient.Database("test").Collection(gamesCollection).UpdateOne(ctx, filter, update, opts)

	if err != nil {
		return UnknownError
	}

	return nil
}
