package node

// This is from: https://gist.github.com/goliatone/e9c13e5f046e34cef6e150d06f20a34c
// Most of it :D

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/ssh"
	"os"
)

func LoadRsaPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

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
	privateDER := x509.MarshalPKCS1PrivateKey(privateKey)
	pemBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&pemBlock)
	return privatePEM
}

func EncodePublicKeyToPEM(publicKey *rsa.PublicKey) []byte {
	publicDER := x509.MarshalPKCS1PublicKey(publicKey)
	pemBlock := pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   publicDER,
	}

	publicPEM := pem.EncodeToMemory(&pemBlock)
	return publicPEM
}

func DecodePublicKeyFromPEM(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
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
