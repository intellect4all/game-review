package reviews

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"
	"sync"
	"time"
)

type Service struct {
	repository Repository
}

type AddReview struct {
	Rating  int    `json:"rating" validate:"required,number,gte=0,lte2"`
	Comment string `json:"comment" validate:"required,min=5,max=2000"`
	GameId  string `json:"gameId" validate:"required"`
	UserId  string `json:"userId" validate:"required"`
}

type ReviewResponse struct {
	Review Review `json:"review"`
	User   User   `json:"user"`
	Vote   Vote   `json:"vote"`
}

func (r *Review) String() string {
	return strings.Join([]string{
		"Review{",
		fmt.Sprintf("Rating: %d", r.Rating),
		"Comment: " + r.Comment,
		"CreatedAt: " + r.CreatedAt.String(),
		"LastUpdatedAt: " + r.LastUpdatedAt.String(),
		"GameId: " + r.GameId,
		"Id: " + r.Id.Hex(),
		fmt.Sprintf("IsDeleted: %t", r.IsDeleted),
		fmt.Sprintf("IsFlagged: %t", r.IsFlagged),
		fmt.Sprintf("Votes: %d", r.Votes),
		"}",
	}, ",")
}

type Vote struct {
	UserId     string `json:"userId" bson:"userId"`
	ReviewId   string `json:"reviewId" bson:"reviewId"`
	IsUpVote   bool   `json:"isUpVote" bson:"isUpVote"`
	IsDownVote bool   `json:"isDownVote" bson:"isDownVote"`
}

type User struct {
	Username string `json:"username" bson:"username"`
	FullName string `json:"fullName" bson:"fullName"`
	Avatar   string `json:"avatar" bson:"displayPic"`
	UserId   string `json:"id" bson:"id"`
}

type Review struct {
	Rating        int                `json:"rating"`
	Comment       string             `json:"comment"`
	CreatedAt     time.Time          `json:"createdAt"`
	LastUpdatedAt time.Time          `json:"lastUpdatedAt"`
	GameId        string             `json:"gameId"`
	Id            primitive.ObjectID `json:"id"`
	IsDeleted     bool               `json:"isDeleted"`
	IsFlagged     bool               `json:"isFlagged"`
	Votes         int                `json:"votes"`
	UserId        string             `json:"userId"`
}

type PaginatedResponseType interface {
	ReviewResponse | Review
}

type PaginatedResponse[V PaginatedResponseType] struct {
	Data         []V  `json:"data"`
	CurrentPage  int  `json:"currentPage"`
	TotalPages   int  `json:"totalPages"`
	TotalItems   int  `json:"totalItems"`
	HasMore      bool `json:"hasMore"`
	ItemsPerPage int  `json:"itemsPerPage"`
}

type Sort struct {
	Key string `json:"sortKey"`
	Asc bool   `json:"asc"`
}

type Repository interface {
	AddReview(ctx context.Context, review *Review) error
	GameExists(ctx context.Context, id string) (bool, error)
	GetReview(ctx context.Context, id string) (*Review, *User, error)
	UpdateReview(ctx context.Context, review *Review) error
	GetReviewsForGame(ctx context.Context, req *GetReviewsForGame) (*PaginatedResponse[ReviewResponse], error)
	GetReviewsForUser(ctx context.Context, req *GetReviewsForGame) (*PaginatedResponse[ReviewResponse], error)
	Vote(ctx context.Context, req VoteRequest, shouldUpvote bool) error
	GetVote(ctx context.Context, userId string, gameId string) (*Vote, error)
	GetFlaggedReviews(ctx context.Context, gameId string, limit int, offset int) (*PaginatedResponse[Review], error)
}

func NewService(r Repository) *Service {
	return &Service{
		repository: r,
	}
}

func (s *Service) addReview(ctx context.Context, r *AddReview) (string, error) {

	exists, err := s.repository.GameExists(ctx, r.GameId)
	if err != nil {
		if !exists {
			return "", ErrGameNotFound
		}

		return "", err
	}

	// create review
	review := getReviewFromAddReview(r)

	// save review
	err = s.repository.AddReview(ctx, &review)

	if err != nil {
		return "", err
	}

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(2)

	//TODO: we can queue this up to be processed later
	go func() {
		defer waitGroup.Done()
		s.checkForPossibleOffensiveContent(ctx, review)
	}()

	go func() {
		defer waitGroup.Done()
		s.updateReviewStats(ctx, r)
	}()

	waitGroup.Wait()

	return review.Id.Hex(), nil
}

type T struct {
	Review
	User
}

func (s *Service) getReview(ctx context.Context, id string) (*ReviewResponse, error) {
	userId := ctx.Value("userId").(string)

	review, user, err := s.repository.GetReview(ctx, id)
	if err != nil {
		return nil, err
	}
	if review == nil || user == nil {
		return nil, ErrReviewNotFound
	}

	if review.IsDeleted || review.IsFlagged {
		return nil, ErrReviewNotFound
	}

	vote, _ := s.repository.GetVote(ctx, userId, id)

	return &ReviewResponse{
		Review: *review,
		User:   *user,
		Vote:   *vote,
	}, nil
}

type GetReviewsForGame struct {
	GameId string `json:"gameId"`
	Limit  int    `json:"limit" validate:"required,number,gte=0,lte=100"`
	Offset int    `json:"offset" validate:"required,number,gte=0"`
	UserId string `json:"userId,omitempty"`
	SortBy Sort   `json:"sortBy,omitempty"`
}

func (s *Service) getReviewsForGame(ctx context.Context, req GetReviewsForGame) (*PaginatedResponse[ReviewResponse], error) {

	reviews, err := s.repository.GetReviewsForGame(ctx, &req)

	if err != nil {
		return nil, err
	}

	if reviews == nil {
		return nil, ErrReviewNotFound
	}

	return reviews, nil
}

func (s *Service) getReviewsForUser(ctx context.Context, req GetReviewsForGame) (*PaginatedResponse[ReviewResponse], error) {
	role := ctx.Value("role").(string)
	userId := ctx.Value("userId").(string)

	if req.UserId == "" || req.UserId != userId || (role != "admin" && role != "moderator") {
		return nil, ErrReviewNotFound
	}

	reviews, err := s.repository.GetReviewsForUser(ctx, &req)

	if err != nil {
		return nil, err
	}

	if reviews == nil {
		return nil, ErrReviewNotFound
	}

	return reviews, nil
}

func (s *Service) updateReview(ctx context.Context, id string, r *AddReview) error {
	userId := ctx.Value("userId").(string)
	role := ctx.Value("role").(string)
	oldReview, user, err := s.repository.GetReview(ctx, id)

	if err != nil {
		return err
	}

	if user == nil || user.UserId != userId || oldReview == nil || (role != "admin" && role != "moderator") {
		return ErrReviewNotFound
	}

	mergeReviews(oldReview, r)

	err = s.repository.UpdateReview(ctx, oldReview)
	if err != nil {
		return err
	}

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(2)

	if strings.TrimSpace(r.Comment) != "" {
		go func() {
			defer waitGroup.Done()
			s.checkForPossibleOffensiveContent(ctx, *oldReview)
		}()
	}

	if r.Rating > 0 {
		go func() {
			defer waitGroup.Done()
			s.updateReviewStats(ctx, r)
		}()
	}

	waitGroup.Wait()

	return nil
}

func (s *Service) deleteReview(ctx context.Context, id string) error {
	userId := ctx.Value("userId").(string)
	role := ctx.Value("role").(string)
	review, u, err := s.repository.GetReview(ctx, id)

	if err != nil {
		return err
	}

	if u == nil || u.UserId != userId || (role != "admin" && role != "moderator") {
		return ErrReviewNotFound
	}

	if review == nil {
		return ErrReviewNotFound
	}

	review.IsDeleted = true

	err = s.repository.UpdateReview(ctx, review)
	if err != nil {
		return err
	}

	go func() {
		//TODO: notify user that their review was deleted
		log.Println("Notifying user that their review was deleted: " + review.String())
	}()

	return nil
}

type VoteRequest struct {
	ReviewId string `json:"reviewId" validate:"required"`
	UpVote   bool   `json:"upVote" validate:"required"`
	UserId   string `json:"userId" validate:"required"`
}

func (s *Service) voteReview(ctx context.Context, reviewId string, shouldUpvote bool) error {
	userId := ctx.Value("userId").(string)

	voteReq := VoteRequest{
		ReviewId: reviewId,
		UpVote:   shouldUpvote,
		UserId:   userId,
	}
	return s.repository.Vote(ctx, voteReq, shouldUpvote)
}

func (s *Service) flagReview(ctx context.Context, id string, flag bool) error {
	role := ctx.Value("role").(string)

	if role != "admin" && role != "moderator" {
		return ErrUnauthorized
	}

	review, _, err := s.repository.GetReview(ctx, id)

	if err != nil {
		return err
	}

	if review == nil {
		return ErrReviewNotFound
	}

	review.IsFlagged = flag

	err = s.repository.UpdateReview(ctx, review)

	if err != nil {
		return err
	}

	go func() {
		//TODO: notify user that their review was deleted
		log.Println("Notifying user that their review was deleted: " + review.String())
	}()

	return nil
}

func (s *Service) getFlaggedReviews(ctx context.Context, gameId string, limit int, offset int) (*PaginatedResponse[Review], error) {
	role := ctx.Value("role").(string)

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	if role != "admin" && role != "moderator" {
		return nil, ErrUnauthorized
	}

	reviews, err := s.repository.GetFlaggedReviews(ctx, gameId, limit, offset)

	if err != nil {
		return nil, err
	}

	if reviews == nil {
		return nil, ErrReviewNotFound
	}

	return reviews, nil
}

func mergeReviews(review *Review, r *AddReview) {
	if strings.TrimSpace(r.Comment) != "" {
		review.Comment = r.Comment
	}

	if r.Rating > 0 {
		review.Rating = r.Rating
	}

	review.LastUpdatedAt = time.Now()
}

func (s *Service) checkForPossibleOffensiveContent(ctx context.Context, review Review) {
	// check for offensive content

}

func (s *Service) updateReviewStats(ctx context.Context, r *AddReview) {
	// update game stats

}

func getReviewFromAddReview(r *AddReview) Review {
	return Review{
		Rating:        r.Rating,
		Comment:       r.Comment,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		GameId:        r.GameId,
		Id:            primitive.NewObjectID(),
		IsDeleted:     false,
		IsFlagged:     false,
		Votes:         0,
		UserId:        r.UserId,
	}

}
