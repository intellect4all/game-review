package reviews

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math"
	"strings"
	"sync"
	"time"
)

type Service struct {
	repository Repository
}

type AddReview struct {
	Rating   int      `json:"rating" validate:"required,number,gte=0,lte2"`
	Comment  string   `json:"comment" validate:"required,min=5,max=2000"`
	GameId   string   `json:"gameId" validate:"required"`
	UserId   string   `json:"userId" validate:"required"`
	Location Location `json:"location" validate:"required"`
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
	Username string   `json:"username" bson:"username"`
	FullName string   `json:"fullName" bson:"fullName"`
	Avatar   string   `json:"avatar" bson:"displayPic"`
	UserId   string   `json:"id" bson:"id"`
	Location Location `json:"location" bson:"location"`
}

type Location struct {
	City        string  `json:"city,omitempty" bson:"city,omitempty"`
	Country     string  `json:"country" bson:"country"`
	Latitude    float64 `json:"latitude" bson:"latitude"`
	Longitude   float64 `json:"longitude" bson:"longitude"`
	CountryCode string  `json:"countryCode" bson:"countryCode"`
}

type Review struct {
	Rating        int                `json:"rating" bson:"rating"`
	Comment       string             `json:"comment" bson:"comment"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	LastUpdatedAt time.Time          `json:"lastUpdatedAt" bson:"lastUpdatedAt"`
	GameId        string             `json:"gameId" bson:"gameId"`
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	IsDeleted     bool               `json:"isDeleted" bson:"isDeleted"`
	IsFlagged     bool               `json:"isFlagged" bson:"isFlagged"`
	Votes         int                `json:"votes"`
	UserId        string             `json:"userId"bson:"userId"`
	Location      Location           `json:"location" bson:"location"`
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
	UpdateReviewStats(ctx context.Context, id string, rating int, ratingCount int) error
	getReviewersForTimeAgo(ctx context.Context, ago time.Time) (*[]Review, error)
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
		err := s.updateReviewStats(ctx, review.GameId, review.Rating, 1)
		if err != nil {
			return
		}
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
			s.updateReviewStats(ctx, r.GameId, r.Rating-oldReview.Rating, 0)
		}()
	}

	waitGroup.Wait()

	return nil
}

func (s *Service) deleteReview(ctx context.Context, id string) error {
	userId := ctx.Value("userId").(string)
	role := ctx.Value("role").(string)
	review, u, err := s.repository.GetReview(ctx, id)

	log.Println("Deleting review: " + id)

	if err != nil {
		return err
	}

	if (u == nil || u.UserId != userId) && (role != "admin" && role != "moderator") {

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
		// update review stats
		err := s.updateReviewStats(ctx, review.GameId, -review.Rating, -1)
		if err != nil {
			log.Println(err)
		}

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
	// check for offensive words
	log.Println("Checking for offensive words in review: " + review.String())
	for _, word := range offensiveWords {
		if strings.Contains(strings.ToLower(review.Comment), word) {
			log.Println("Found offensive word: " + word)
			go func(reviewId string) {
				err := s.flagReview(ctx, reviewId, true)
				if err != nil {
					log.Println("Error flagging review: " + err.Error())
					return
				}
			}(review.Id.Hex())
			return
		}
	}

}

func (s *Service) updateReviewStats(ctx context.Context, gameId string, rating int, ratingCount int) error {
	// update game stats
	log.Println("Updating review stats for game: " + gameId)
	return s.repository.UpdateReviewStats(ctx, gameId, rating, ratingCount)

}

type LocationReqType int

const (
	Day LocationReqType = iota
	Week
	Month
)

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type UserWithCount struct {
	UserId   string   `json:"id" bson:"_id"`
	Count    int      `json:"count" bson:"count"`
	Location Location `json:"location" bson:"location"`
}

func (s *Service) getLocation(ctx context.Context, reqType LocationReqType, value int) ([]LocationWithCount, error) {
	log.Println("inside 1")
	timeAgo := getTimeAgoFromReqType(reqType, value)

	reviews, err := s.repository.getReviewersForTimeAgo(ctx, timeAgo)

	log.Println("inside 2")

	if err != nil {
		return nil, err
	}

	reviewLocations, err := sortUserLocationProximity(reviews)

	log.Println("inside 3")

	if err != nil {
		return nil, err
	}

	return *reviewLocations, nil
}

type LocationWithCount struct {
	Location LatLng `json:"location"`
	Count    int    `json:"count"`
}

func sortUserLocationProximity(reviews *[]Review) (*[]LocationWithCount, error) {
	log.Println("sort user location proximity")
	proximityInKm := 0.1

	locationsWithCount := make([]LocationWithCount, 0)

	checkedIndexesMap := make(map[int]bool)

	for i, review := range *reviews {

		log.Println("checking ", i)

		if checkedIndexesMap[i] {
			log.Println("already checked")
			continue
		}

		latLng := LatLng{
			Lat: review.Location.Latitude,
			Lng: review.Location.Longitude,
		}

		currentLoc := LocationWithCount{
			Location: latLng,
			Count:    1,
		}

		// print current location
		log.Println("current location: ", currentLoc)

		for j := i + 1; j < len(*reviews); j++ {

			log.Println("checking ", j)

			if checkedIndexesMap[j] {
				log.Println(j, " already checked")
				continue
			}

			thisLoc := LocationWithCount{
				Location: LatLng{
					Lat: (*reviews)[j].Location.Latitude,
					Lng: (*reviews)[j].Location.Longitude,
				},
				Count: 1,
			}

			if isLocationClose(currentLoc.Location, thisLoc.Location, proximityInKm) {
				currentLoc.Count++
				checkedIndexesMap[j] = true
			}
		}

		locationsWithCount = append(locationsWithCount, currentLoc)
	}

	log.Println("sort user location proximity 2")
	return &locationsWithCount, nil
}

func distance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	R := 6371000.0 // metres

	a := 0.5 - math.Cos((lat2-lat1)*math.Pi/180)/2 + math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*(1-math.Cos((lon2-lon1)*math.Pi/180))/2

	return R * 2.0 * math.Asin(math.Sqrt(a))

}

func isLocationClose(location1 LatLng, location2 LatLng, proximityInKm float64) bool {
	log.Println("is location close")

	// check if location is close
	distance := distance(location1.Lat, location1.Lng, location2.Lat, location2.Lng)

	log.Println("distance: ", distance)

	return distance < float64(proximityInKm*1000)
}

func getTimeAgoFromReqType(reqType LocationReqType, value int) time.Time {
	if value == 0 {
		value = 1
	}
	switch reqType {
	case Day:
		return time.Now().AddDate(0, 0, -value)
	case Week:
		return time.Now().AddDate(0, 0, -value*7)
	case Month:
		return time.Now().AddDate(0, -value, 0)
	default:
		return time.Now().AddDate(0, 0, -value)
	}
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

var offensiveWords []string

func init() {
	offensiveWords = []string{
		"fuck",
		"shit",
		"fuck off",
		"cunt",
		"motherfucker",
		"fucker",
		"wanker",
		"dumbass",
		"nigga",
		"nigger",
	}
}
