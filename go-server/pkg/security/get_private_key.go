package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GetPrivateKey(file string) (*rsa.PrivateKey, error) {
	secretKey, err := GetPrivateKeyBytes(file)

	if err != nil {
		return nil, err
	}

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
			return
		}

	}(secretPrivateFile)

	if err != nil {
		return nil, err
	}

	fileStat, err := secretPrivateFile.Stat()

	if err != nil {
		return nil, err
	}

	secretKey := make([]byte, fileStat.Size())
	_, err = secretPrivateFile.Read(secretKey)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(secretKey)

	if block == nil {
		return nil, err
	}

	return block.Bytes, nil
}
