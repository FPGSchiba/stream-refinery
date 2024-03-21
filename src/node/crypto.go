package node

// This is from: https://gist.github.com/goliatone/e9c13e5f046e34cef6e150d06f20a34c

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"os"
)

func LoadRsaPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParseRawPrivateKey(data)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}

func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privateDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	pemBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&pemBlock)

	return privatePEM
}

func GeneratePublicKey(privateKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}

func WriteKeyToFile(keyBytes []byte, saveFileTo string) error {
	file, err := os.OpenFile(saveFileTo, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	_, err = file.Write(keyBytes)
	if err != nil {
		return err
	}

	return nil
}
