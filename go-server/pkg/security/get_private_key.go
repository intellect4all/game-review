package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GetPrivateKey(file string) (*rsa.PrivateKey, error) {
	secretKey, err := GetPrivateKeyBytes(file)

	privateKey, err := x509.ParsePKCS1PrivateKey(secretKey)
	if err != nil {
		return nil, err
	}

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

	secretKey := make([]byte, fileStat.Size())
	_, err = secretPrivateFile.Read(secretKey)

	if err != nil {
		fmt.Printf("Error reading key: %s", err)
		return nil, err
	}

	block, _ := pem.Decode(secretKey)

	return block.Bytes, nil
}
