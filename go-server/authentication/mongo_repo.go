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
	"unicode"
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

	valRes := m.validate.Struct(userCredential)

	if valRes != nil {
		return valRes
	}

	// check if user already exists
	var userCredentialFromDb UserCredential

	err := m.mongoDbClient.Database("test").Collection("users").FindOne(ctx, bson.D{{"email", userCredential.Id}}).Decode(&userCredentialFromDb)

	// check if the error is a ErrNoDocuments error
	if err == nil {
		return ErrUserAlreadyExists
	}

	if ok, errMessage := isPasswordValid(userCredential.Password); !ok {
		return errors.New("Invalid Password: " + errMessage)
	}

	// encrypt password
	hashedPassword, err := encryptPassword(userCredential.Password)
	if err != nil {
		return err
	}

	userCredential.Password = hashedPassword

	_, err = m.mongoDbClient.Database("test").Collection("users").InsertOne(ctx, userCredential)
	if err != nil {
		return err
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

func (m *MongoRepository) isValidPassword(password string, encryptedPassword string) bool {
	return isCorrectPassword(password, encryptedPassword)
}

func (m *MongoRepository) CreateVerificationOTP(ctx context.Context, id *UserID) (string, error) {
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

func (m *MongoRepository) VerifyUser(ctx context.Context, requestData *VerifyUserRequest) error {
	userCred, err := m.GetUserCredential(ctx, UserID(requestData.Email))
	if err != nil {
		return ErrUserNotFound
	}
	// convert tokenID to objectID
	id, _ := primitive.ObjectIDFromHex(requestData.TokenID)

	// check if otp code is valid

	var otpData OtpData
	res := m.mongoDbClient.Database("test").Collection("otpCodes").FindOne(ctx, bson.M{"_id": id}).Decode(&otpData)

	if res != nil {
		return ErrInvalidOTP
	}

	if otpData.Used {
		return ErrOTPUsed
	}

	if otpData.ExpirationTime.Before(time.Now()) {
		return ErrOTPExpired
	}

	if otpData.OtpCode != requestData.OTPCode {
		return ErrInvalidOTP
	}

	// update otp code to used
	_, err = m.mongoDbClient.Database("test").Collection("otpCodes").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"used": true}})

	if err != nil {
		return err
	}

	// update user to verified
	_, err = m.mongoDbClient.Database("test").Collection("users").UpdateOne(ctx, bson.M{"email": userCred.Id}, bson.M{"$set": bson.M{"isVerified": true}})

	if err != nil {
		return UnknownError
	}

	return nil
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

func isCorrectPassword(password string, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func isPasswordValid(password string) (bool, string) {
	var (
		upp, low, num, sym bool
		tot                uint8
		errorMessage       string
	)

	errorMessage = "Password must contain at least 8 characters, 1 uppercase, 1 lowercase, 1 number, and 1 symbol."

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			errorMessage = "Password contains invalid characters. Only letters, numbers, and symbols are allowed."
		}
	}

	if !upp || !low || !num || !sym || tot < 8 {

		return false, errorMessage
	}

	return true, ""
}
