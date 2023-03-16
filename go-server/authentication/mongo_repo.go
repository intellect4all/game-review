package authentication

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type MongoRepository struct {
	mongoDbClient *mongo.Client
	validate      *validator.Validate
}

func NewMongoRepository(mongoDbClient *mongo.Client) *MongoRepository {
	return &MongoRepository{
		mongoDbClient: mongoDbClient,
		validate:      validator.New(),
	}
}

func (m *MongoRepository) CreateNewUser(ctx context.Context, userCredential *UserCredential) error {

	_, err := m.mongoDbClient.Database("test").Collection("users").InsertOne(ctx, userCredential)

	if err != nil {
		return UnknownError
	}

	return nil
}

func (m *MongoRepository) GetUserCredential(ctx context.Context, userId UserID) (*UserCredential, error) {
	var userCredential UserCredential

	fmt.Printf("userIDString %s \n", userId)
	err := m.mongoDbClient.Database("test").Collection("users").FindOne(ctx, bson.D{{"email", userId}}).Decode(&userCredential)
	if err != nil {
		fmt.Printf("Error: %v, %v \n", err, err.Error())
		return nil, ErrUserNotFound
	}

	return &userCredential, nil
}

func (m *MongoRepository) GetUserDetail(ctx context.Context, userId UserID) (*UserDetail, error) {
	res := m.mongoDbClient.Database("test").Collection("users").FindOne(ctx, bson.D{{"_id", string(userId)}})

	if res.Err() != nil {
		return nil, res.Err()
	}

	var userDetail *UserDetail
	err := res.Decode(&userDetail)
	if err != nil {
		return nil, err
	}

	return userDetail, nil
}

func (m *MongoRepository) DeleteUser(ctx context.Context, userId UserID) error {
	res, err := m.mongoDbClient.Database("test").Collection("users").UpdateOne(ctx, bson.D{{"_id", string(userId)}}, bson.D{{"$set", bson.D{{"deleted", true}}}})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

//func (m *MongoRepository) isCorrectPassword(password string, encryptedPassword string) bool {
//	return isCorrectPassword(password, encryptedPassword)
//}

func (m *MongoRepository) CreateOTP(ctx context.Context, id *UserID) (string, error) {
	otpData := getOTPData(id)

	fmt.Printf("otpData %v \n", otpData)
	res, err := m.mongoDbClient.Database("test").Collection("otpCodes").InsertOne(ctx, otpData)
	if err != nil {
		return "", ErrOTPCreationFailed
	}

	var idd primitive.ObjectID
	idd = res.InsertedID.(primitive.ObjectID)

	otpID := idd.Hex()

	return otpID, nil
}

func (m *MongoRepository) VerifyUser(ctx context.Context, requestData *OtpCheckData) error {
	userCred, err := m.GetUserCredential(ctx, UserID(requestData.Email))
	if err != nil {
		return ErrUserNotFound
	}

	// check if otp code is valid
	if _, err := m.VerifyOTP(ctx, requestData); err != nil {
		return err
	}

	// update user to verified
	_, err = m.mongoDbClient.Database("test").Collection("users").UpdateOne(ctx, bson.M{"email": userCred.Id}, bson.M{"$set": bson.M{"isVerified": true}})

	if err != nil {
		return UnknownError
	}

	return nil
}

func (m *MongoRepository) ChangePassword(ctx context.Context, f *ForgetAndResetPasswordRequest) error {
	_, err := m.GetUserCredential(ctx, UserID(f.Email))
	if err != nil {
		return ErrUserNotFound
	}

	// check if otp code is valid
	otpData := OtpCheckData{
		Email:   f.Email,
		TokenID: f.TokenID,
		OTPCode: f.OTPCode,
	}
	if _, err := m.VerifyOTP(ctx, &otpData); err != nil {
		return err
	}

	if ok, errMessage := isPasswordValid(f.Password); !ok {
		return errors.New("Invalid Password: " + errMessage)
	}

	if f.Password != f.ConfirmPassword {
		return ErrPasswordMismatch
	}

	// encrypt password
	hashedPassword, err := encryptPassword(f.Password)
	if err != nil {
		return err
	}

	// update user password
	_, err = m.mongoDbClient.Database("test").Collection("users").UpdateOne(ctx, bson.M{"email": f.Email}, bson.M{"$set": bson.M{"password": hashedPassword}})

	if err != nil {
		return UnknownError
	}

	return nil
}

func (m *MongoRepository) VerifyOTP(ctx context.Context, requestData *OtpCheckData) (bool, error) {
	// convert tokenID to objectID
	id, _ := primitive.ObjectIDFromHex(requestData.TokenID)

	var otpData OtpData
	res := m.mongoDbClient.Database("test").Collection("otpCodes").FindOne(ctx, bson.M{"_id": id}).Decode(&otpData)

	if res != nil {
		return false, ErrInvalidOTP
	}

	if otpData.Used {
		return false, ErrOTPUsed
	}

	if otpData.ExpirationTime.Before(time.Now()) {
		return false, ErrOTPExpired
	}

	if otpData.OtpCode != requestData.OTPCode {
		return false, ErrInvalidOTP
	}

	// update otp code to used
	_, err := m.mongoDbClient.Database("test").Collection("otpCodes").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"used": true}})

	if err != nil {
		return false, UnknownError
	}

	return true, nil
}

type OtpData struct {
	UserId         string    `bson:"userId"`
	OtpCode        string    `bson:"otpCode"`
	Used           bool      `bson:"used"`
	CreatedTime    time.Time `bson:"createdAt"`
	ExpirationTime time.Time `bson:"expiresAt"`
}

func getOTPData(u *UserID) *OtpData {
	otpCode := generateOTPCode()

	otpData := OtpData{
		UserId:         string(*u),
		OtpCode:        otpCode,
		Used:           false,
		CreatedTime:    time.Now(),
		ExpirationTime: time.Now().Add(5 * time.Minute),
	}

	return &otpData

}

func generateOTPCode() string {
	// generate random 6 digit otp code
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := r.Intn(999999-100000) + 100000
	return fmt.Sprintf("%d", code)
}

func encryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
