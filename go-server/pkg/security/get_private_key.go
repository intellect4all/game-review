package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func GetPrivateKey(file string) (*rsa.PrivateKey, error) {
	log.Println("Getting private key")
	secretKey, err := GetPrivateKeyBytes(file)

	if err != nil {
		log.Println("Error getting private key")
		return nil, err
	}

	log.Println("Private key obtained")

	privateKey, err := x509.ParsePKCS1PrivateKey(secretKey)
	if err != nil {
		fmt.Printf("Error parsing private key: %s", err.Error())
		return nil, err
	}

	log.Println("Private key parsed")

	return privateKey, nil
}

func GetPrivateKeyBytes(file string) ([]byte, error) {
	file = "private/" + file
	secretPrivateFile, err := os.Open(file)
	defer func(secretPrivateFile *os.File) {
		err := secretPrivateFile.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s", err)
			return
		}
		fmt.Printf("Closed file: %s", file)
	}(secretPrivateFile)

	if err != nil {
		fmt.Printf("Error opening file: %s", err)

		return nil, err
	}

	fileStat, err := secretPrivateFile.Stat()

	if err != nil {
		fmt.Printf("Error getting stat file: %s", err)
		return nil, err
	}

	log.Println("Pirvate stat obtained")

	secretKey := make([]byte, fileStat.Size())
	_, err = secretPrivateFile.Read(secretKey)

	log.Println("byte read obtained " + string(secretKey) + " " + file)

	if err != nil {
		fmt.Printf("Error reading key: %s", err)
		return nil, err
	}

	block, rest := pem.Decode(secretKey)

	if block == nil {
		fmt.Printf("Error decoding key, %s", rest)
		return nil, err
	}

	log.Println("Private key bytes obtained " + string(block.Bytes) + " " + file)
	return block.Bytes, nil
}

func GetPublicKey(file string) ([]byte, error) {
	file = "private/" + file
	secretPrivateFile, err := os.Open(file)
	defer func(secretPrivateFile *os.File) {
		err := secretPrivateFile.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s", err)
			return
		}
		fmt.Printf("Closed file: %s", file)
	}(secretPrivateFile)

	if err != nil {
		fmt.Printf("Error opening file: %s", err)

		return nil, err
	}

	fileStat, err := secretPrivateFile.Stat()

	if err != nil {
		fmt.Printf("Error getting stat file: %s", err)
		return nil, err
	}

	log.Println("Pirvate stat obtained")

	secretKey := make([]byte, fileStat.Size())
	_, err = secretPrivateFile.Read(secretKey)

	log.Println("byte read obtained " + string(secretKey) + " " + file)

	if err != nil {
		fmt.Printf("Error reading key: %s", err)
		return nil, err
	}

	return secretKey, nil
}
