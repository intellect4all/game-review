package notifications

import (
	"crypto/tls"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"os"
	"strconv"
	"time"
)

type EmailRequest struct {
	To      string
	From    string
	Subject string
	Body    string
}

type SmtpDetails struct {
	Host     string
	Port     int
	Username string
	Password string
}

func SendEmail(request EmailRequest) error {
	server := mail.NewSMTPClient()

	// SMTP Server

	smtpDetails := getSmtpDetails()
	server.Host = smtpDetails.Host
	server.Port = smtpDetails.Port
	server.Username = smtpDetails.Username
	server.Password = smtpDetails.Password
	server.Encryption = mail.EncryptionTLS

	// Variable to keep alive connection
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	log.Println("Connecting to SMTP Server")

	// SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Connected to SMTP Server")

	// New email simple html with inline and CC
	email := mail.NewMSG()
	email.SetFrom(request.From).
		AddTo(request.To).
		SetSubject(request.Subject)

	email.SetBody(mail.TextHTML, request.Body)

	//TODO: get private key

	//privateKey, err := security.GetPrivateKeyBytes("key.pem")

	//if err != nil {
	//	log.Println(err)
	//}

	// you can add dkim signature to the email.
	// to add dkim, you need a private key already created one.
	//if len(privateKey) > 0 {
	//	options := dkim.NewSigOptions()
	//	options.PrivateKey = privateKey
	//	options.Domain = "cool_game_review.com"
	//	options.Selector = "default"
	//	options.SignatureExpireIn = 3600
	//	options.Headers = []string{"from", "date", "mime-version", "received", "received"}
	//	options.AddSignatureTimestamp = true
	//	options.Canonicalization = "relaxed/relaxed"
	//
	//	email.SetDkim(options)
	//}

	// always check error after send
	if email.Error != nil {
		log.Fatal(email.Error)
	}

	// Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent")
	}

	return err
}

func getSmtpDetails() *SmtpDetails {
	var details SmtpDetails

	details.Host = os.Getenv("SMTP_HOST")
	details.Port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	details.Username = os.Getenv("SMTP_USERNAME")
	details.Password = os.Getenv("SMTP_PASSWORD")

	return &details
}
