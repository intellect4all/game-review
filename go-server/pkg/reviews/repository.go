package reviews

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"sync"
)

type RepositoryImpl struct {
	mongoDbClient *mongo.Client
	validate      *validator.Validate
}

const (
	reviewsCollection = "reviews"
)

func NewRepository(mongoClient *mongo.Client) *RepositoryImpl {
	return &RepositoryImpl{
		mongoDbClient: mongoClient,
	}
}

func (r *RepositoryImpl) AddReview(ctx context.Context, review *Review) error {
	_, err := r.mongoDbClient.Database("test").Collection(reviewsCollection).InsertOne(ctx, review)

	if err != nil {
		log.Println(err)
		return UnknownError
	}

	return nil
}

func (r *RepositoryImpl) UpdateReview(ctx context.Context, review *Review) error {

	filter := bson.D{{"_id", review.Id}}
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", review}}

	_, err := r.mongoDbClient.Database("test").Collection(reviewsCollection).UpdateOne(ctx, filter, update, opts)

	if err != nil {
		log.Println(err)
		return UnknownError
	}

	return nil
}

type UserRw struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Avatar    string             `bson:"displayPic"`
	Username  string             `bson:"username"`
	firstName string             `bson:"firstName"`
	lastName  string             `bson:"lastName"`
}

func (r *RepositoryImpl) GetReview(ctx context.Context, id string) (*Review, *User, error) {

	var review Review
	var userRw UserRw

	filter := bson.D{{"_id", id}, {"isDeleted", false}}

	err := r.mongoDbClient.Database("test").Collection(reviewsCollection).FindOne(ctx, filter).Decode(&review)

	if err != nil {
		log.Println(err)
		return nil, nil, ErrNotFound
	}

	filter = bson.D{{"_id", review.UserId}, {"isDeleted", false}}
	returnedFields := bson.D{{"_id", 1}, {"displayPic", 1}, {"username", 1}, {"firstName", 1}, {"lastName", 1}}
	opts := options.FindOne().SetProjection(returnedFields)

	err = r.mongoDbClient.Database("test").Collection("users").FindOne(ctx, filter, opts).Decode(&userRw)
	if err != nil {
		log.Println(err)
		return nil, nil, ErrNotFound
	}

	var user User

	user.UserId = userRw.Id.Hex()
	user.Avatar = userRw.Avatar
	user.Username = userRw.Username
	user.FullName = userRw.firstName + " " + userRw.lastName

	return &review, &user, nil
}

func (r *RepositoryImpl) GameExists(ctx context.Context, id string) (bool, error) {
	res := r.mongoDbClient.Database("test").Collection("games").FindOne(ctx, bson.D{{"_id", id}, {"isDeleted", false}})

	if res.Err() != nil {
		return false, nil
	}

	return true, nil
}

func (r *RepositoryImpl) GetVote(ctx context.Context, userId string, reviewId string) (*Vote, error) {
	var vote Vote

	defaultVote := &Vote{
		UserId:     userId,
		ReviewId:   reviewId,
		IsUpVote:   false,
		IsDownVote: false,
	}

	filter := bson.D{{"userId", userId}, {"reviewId", reviewId}}
	res := r.mongoDbClient.Database("test").Collection("votes").FindOne(ctx, filter)

	if res.Err() != nil {
		return defaultVote, nil
	}

	err := res.Decode(&vote)

	if err != nil {
		log.Println(err)
		return defaultVote, UnknownError
	}

	return &vote, nil

}

func (r *RepositoryImpl) Vote(ctx context.Context, req VoteRequest, shouldUpVote bool) error {
	// increment or decrement the vote count
	filter := bson.D{{"_id", req.ReviewId}}

	voteFilter := bson.D{{"userId", req.UserId}, {"reviewId", req.ReviewId}}

	opts := options.Update().SetUpsert(true)

	// check if the user has already voted
	voteRes := r.mongoDbClient.Database("test").Collection("votes").FindOne(ctx, voteFilter)

	if voteRes.Err() != nil {
		// user has not voted
		// add a new vote
		_, err := r.mongoDbClient.Database("test").Collection("votes").InsertOne(ctx, Vote{
			UserId:     req.UserId,
			ReviewId:   req.ReviewId,
			IsUpVote:   shouldUpVote,
			IsDownVote: !shouldUpVote,
		})

		if err != nil {
			log.Println(err)
			return UnknownError
		}

		// increment the vote count

		inc := 0

		if shouldUpVote {
			inc = 1
		} else {
			inc = -1
		}

		update := bson.D{{"$inc", bson.D{{"votes", inc}}}}

		_, err = r.mongoDbClient.Database("test").Collection(reviewsCollection).UpdateOne(ctx, filter, update, opts)

		if err != nil {
			log.Println(err)
			return UnknownError
		}

		return nil

	} else {
		// user has already voted
		// update the vote
		var vote Vote
		err := voteRes.Decode(&vote)
		if err != nil {
			log.Println(err)
			return UnknownError
		}

		if vote.IsUpVote == shouldUpVote {
			// user is trying to upvote an already upvoted review
			// or
			// user is trying to downvote an already downvoted review
			return nil
		}

		// update the vote
		update := bson.D{{"$set", bson.D{{"isUpVote", shouldUpVote}, {"isDownVote", !shouldUpVote}}}}
		_, err = r.mongoDbClient.Database("test").Collection("votes").UpdateOne(ctx, voteFilter, update)
		if err != nil {
			log.Println(err)
			return UnknownError
		}

		// update the vote count
		inc := 0

		if shouldUpVote {
			inc = 2
		} else {
			inc = -2
		}

		update = bson.D{{"$inc", bson.D{{"votes", inc}}}}

		_, err = r.mongoDbClient.Database("test").Collection(reviewsCollection).UpdateOne(ctx, filter, update, opts)

		if err != nil {
			log.Println(err)
			return UnknownError
		}

		return nil
	}
}

func (r *RepositoryImpl) GetReviewsForGame(ctx context.Context, req *GetReviewsForGame) (*PaginatedResponse[ReviewResponse], error) {
	var reviews []Review
	var userRws []UserRw

	gameRevFilter := bson.D{{"gameId", req.GameId}, {"isDeleted", false}}

	sortVal := -1
	if req.SortBy.Asc {
		sortVal = 1
	}

	opts := options.Find().SetLimit(int64(req.Limit)).SetSkip(int64(req.Offset)).SetSort(bson.D{{req.SortBy.Key, sortVal}})

	cursor, err := r.mongoDbClient.Database("test").Collection(reviewsCollection).Find(ctx, gameRevFilter, opts)

	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	err = cursor.All(ctx, &reviews)
	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	var userIds []string

	for _, review := range reviews {
		userIds = append(userIds, review.UserId)
	}

	userIdFilter := bson.D{{"_id", bson.D{{"$in", req.UserId}}}}

	returnedFields := bson.D{{"_id", 1}, {"displayPic", 1}, {"username", 1}, {"firstName", 1}, {"lastName", 1}}
	opts = options.Find().SetProjection(returnedFields)

	cursor, err = r.mongoDbClient.Database("test").Collection("users").Find(ctx, userIdFilter, opts)

	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	err = cursor.All(ctx, &userRws)
	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	var reviewResponses []ReviewResponse

	for _, review := range reviews {
		var user User
		for _, userRw := range userRws {
			if userRw.Id.Hex() == review.UserId {
				user.UserId = userRw.Id.Hex()
				user.Avatar = userRw.Avatar
				user.Username = userRw.Username
				user.FullName = userRw.firstName + " " + userRw.lastName
			}
		}

		reviewResponses = append(reviewResponses, ReviewResponse{
			Review: review,
			User:   user,
		})
	}

	// get votes for each reviews with go routines
	var wg sync.WaitGroup
	wg.Add(len(reviewResponses))

	for i, reviewResponse := range reviewResponses {
		go func(i int, reviewResponse ReviewResponse) {
			defer wg.Done()
			vote, err := r.GetVote(ctx, req.UserId, reviewResponse.Review.Id.Hex())
			if err != nil {
				log.Println(err)
			}
			reviewResponses[i].Vote = *vote
		}(i, reviewResponse)
	}

	wg.Wait()

	var count int64

	count, err = r.mongoDbClient.Database("test").Collection(reviewsCollection).CountDocuments(ctx, gameRevFilter)

	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	var response *PaginatedResponse[ReviewResponse]

	response = &PaginatedResponse[ReviewResponse]{
		Data:         reviewResponses,
		TotalPages:   int(math.Ceil(float64(count) / float64(req.Limit))),
		CurrentPage:  int(math.Ceil(float64(req.Offset) / float64(req.Limit))),
		TotalItems:   int(count),
		HasMore:      int(count) > (req.Offset + req.Limit),
		ItemsPerPage: req.Limit,
	}

	return response, nil
}

func (r *RepositoryImpl) GetReviewsForUser(ctx context.Context, req *GetReviewsForGame) (*PaginatedResponse[ReviewResponse], error) {
	var reviews []Review
	var rawUser UserRw

	gameRevFilter := bson.D{{"userId", req.UserId}, {"isDeleted", false}}

	sortVal := -1
	if req.SortBy.Asc {
		sortVal = 1
	}

	opts := options.Find().SetLimit(int64(req.Limit)).SetSkip(int64(req.Offset)).SetSort(bson.D{{req.SortBy.Key, sortVal}})

	cursor, err := r.mongoDbClient.Database("test").Collection(reviewsCollection).Find(ctx, gameRevFilter, opts)

	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	err = cursor.All(ctx, &reviews)
	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	var userIds []string

	for _, review := range reviews {
		userIds = append(userIds, review.UserId)
	}

	userIdFilter := bson.D{{"_id", req.UserId}}

	returnedFields := bson.D{{"_id", 1}, {"displayPic", 1}, {"username", 1}, {"firstName", 1}, {"lastName", 1}}
	getUserOpt := options.FindOne().SetProjection(returnedFields)

	err = r.mongoDbClient.Database("test").Collection("users").FindOne(ctx, userIdFilter, getUserOpt).Decode(&rawUser)

	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	err = cursor.All(ctx, &rawUser)
	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	var reviewResponses []ReviewResponse

	user := User{
		UserId:   rawUser.Id.Hex(),
		Avatar:   rawUser.Avatar,
		Username: rawUser.Username,
		FullName: rawUser.firstName + " " + rawUser.lastName,
	}

	for _, review := range reviews {
		reviewResponses = append(reviewResponses, ReviewResponse{
			Review: review,
			User:   user,
		})
	}

	// get votes for each reviews with go routines
	var wg sync.WaitGroup
	wg.Add(len(reviewResponses))

	for i, reviewResponse := range reviewResponses {
		go func(i int, reviewResponse ReviewResponse) {
			defer wg.Done()
			vote, err := r.GetVote(ctx, req.UserId, reviewResponse.Review.Id.Hex())
			if err != nil {
				log.Println(err)
			}
			reviewResponses[i].Vote = *vote
		}(i, reviewResponse)
	}

	wg.Wait()

	var count int64

	count, err = r.mongoDbClient.Database("test").Collection(reviewsCollection).CountDocuments(ctx, gameRevFilter)

	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	var response *PaginatedResponse[ReviewResponse]

	response = &PaginatedResponse[ReviewResponse]{
		Data:         reviewResponses,
		TotalPages:   int(math.Ceil(float64(count) / float64(req.Limit))),
		CurrentPage:  int(math.Ceil(float64(req.Offset) / float64(req.Limit))),
		TotalItems:   int(count),
		HasMore:      int(count) > (req.Offset + req.Limit),
		ItemsPerPage: req.Limit,
	}

	return response, nil
}

func (r *RepositoryImpl) GetFlaggedReviews(ctx context.Context, gameId string, limit int, offset int) (*PaginatedResponse[Review], error) {
	var reviews []Review

	gameRevFilter := bson.D{{"isDeleted", false}, {"isFlagged", true}}

	if gameId != "" {
		gameRevFilter = append(gameRevFilter, bson.E{Key: "gameId", Value: gameId})
	}

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(bson.D{{"createdAt", -1}})

	cursor, err := r.mongoDbClient.Database("test").Collection(reviewsCollection).Find(ctx, gameRevFilter, opts)

	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	err = cursor.All(ctx, &reviews)
	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	var count int64

	count, err = r.mongoDbClient.Database("test").Collection(reviewsCollection).CountDocuments(ctx, gameRevFilter)

	if err != nil {
		log.Println(err)
		return nil, UnknownError
	}

	var response *PaginatedResponse[Review]

	response = &PaginatedResponse[Review]{
		Data:         reviews,
		TotalPages:   int(math.Ceil(float64(count) / float64(limit))),
		CurrentPage:  int(math.Ceil(float64(offset) / float64(limit))),
		TotalItems:   int(count),
		HasMore:      int(count) > (offset + limit),
		ItemsPerPage: limit,
	}

	return response, nil
}
