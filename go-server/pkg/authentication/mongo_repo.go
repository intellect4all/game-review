package authentication

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go-server/pkg/notifications"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
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

func (m *MongoRepository) CreateNewUser(ctx context.Context, userCredential UserCredential) error {

	_, err := m.mongoDbClient.Database("test").Collection("users").InsertOne(ctx, userCredential)

	if err != nil {
		log.Println(err)
		return UnknownError
	}

	return nil
}

func (m *MongoRepository) GetUserCredentialByEmail(ctx context.Context, email string) (*UserCredential, error) {
	return m.getUserCredential(ctx, "email", email)
}

func (m *MongoRepository) GetUserCredentialById(ctx context.Context, id string) (*UserCredential, error) {
	return m.getUserCredential(ctx, "_id", id)
}

func (m *MongoRepository) GetUserCredentialByUserName(ctx context.Context, username string) (*UserCredential, error) {
	var userCredential UserCredential

	err := m.mongoDbClient.Database("test").Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&userCredential)
	if err != nil {

		return nil, ErrUserNotFound
	}
	return &userCredential, nil
}

func (m *MongoRepository) getUserCredential(ctx context.Context, field string, value any) (*UserCredential, error) {
	var userCredential UserCredential

	err := m.mongoDbClient.Database("test").Collection("users").FindOne(ctx, bson.D{{field, value}}).Decode(&userCredential)
	if err != nil {

		return nil, ErrUserNotFound
	}

	return &userCredential, nil
}

func (m *MongoRepository) GetUserDetail(ctx context.Context, email string) (*UserDetail, error) {
	res := m.mongoDbClient.Database("test").Collection("users").FindOne(ctx, bson.D{{"_id", string(email)}})

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

func (m *MongoRepository) DeleteUser(ctx context.Context, email string) error {
	res, err := m.mongoDbClient.Database("test").Collection("users").UpdateOne(ctx, bson.D{{"_id", string(email)}}, bson.D{{"$set", bson.D{{"deleted", true}}}})
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

func (m *MongoRepository) CreateOTP(ctx context.Context, id string) (string, error) {
	otpData := getOTPData(id)

	res, err := m.mongoDbClient.Database("test").Collection("otpCodes").InsertOne(ctx, otpData)
	if err != nil {
		return "", ErrOTPCreationFailed
	}

	var idd primitive.ObjectID
	idd = res.InsertedID.(primitive.ObjectID)

	otpID := idd.Hex()

	sendEmailToUser(otpData)

	return otpID, nil
}

func sendEmailToUser(data *OtpData) {
	emailRequest := notifications.EmailRequest{
		From:    "cool_game_rev.com",
		To:      data.email,
		Subject: "Your OTP Code",
		Body:    GetHtmlTemplate(data),
	}
	err := notifications.SendEmail(
		emailRequest,
	)
	if err != nil {
		return
	}
}

func (m *MongoRepository) VerifyUser(ctx context.Context, requestData VerifyAccountRequest) error {
	userCred, err := m.GetUserCredentialByEmail(ctx, requestData.Email)
	if err != nil {
		return ErrUserNotFound
	}

	// check if otp code is valid
	if _, err := m.VerifyOTP(ctx, &requestData); err != nil {
		return err
	}

	// update user to verified
	_, err = m.mongoDbClient.Database("test").Collection("users").UpdateOne(ctx, bson.M{"email": userCred.Id}, bson.M{"$set": bson.M{"isVerified": true}})

	if err != nil {
		return UnknownError
	}

	return nil
}

func (m *MongoRepository) ChangePassword(ctx context.Context, f ForgetAndResetPasswordRequest) error {
	_, err := m.GetUserCredentialByEmail(ctx, f.Email)
	if err != nil {
		return ErrUserNotFound
	}

	// check if otp code is valid
	otpData := VerifyAccountRequest{
		Email:   f.Email,
		TokenID: f.TokenId,
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

func (m *MongoRepository) VerifyOTP(ctx context.Context, requestData *VerifyAccountRequest) (bool, error) {
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
	email          string    `bson:"email"`
	OtpCode        string    `bson:"otpCode"`
	Used           bool      `bson:"used"`
	CreatedTime    time.Time `bson:"createdAt"`
	ExpirationTime time.Time `bson:"expiresAt"`
}

func getOTPData(u string) *OtpData {
	otpCode := generateOTPCode()

	otpData := OtpData{
		email:          u,
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
		log.Println("Error encrypting password: ", err)
		return "", err
	}

	return string(hashedPassword), nil
}

type OTPHtmlTemplate struct {
	OTPCode   string
	BrandName string
	Address   string
	State     string
}

func GetHtmlTemplate(data *OtpData) string {

	tmplt, _ := template.ParseFiles("resources/templates/otp.html")

	tmplData := OTPHtmlTemplate{
		OTPCode:   data.OtpCode,
		BrandName: "Cool Game",
		Address:   "1234 Main St",
		State:     "CA",
	}

	var tpl bytes.Buffer

	err := tmplt.Execute(&tpl, tmplData)

	if err != nil {
		return ""
	}

	return tpl.String()

}
