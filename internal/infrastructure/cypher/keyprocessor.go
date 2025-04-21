package cypher

import (
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"errors"

	domainErrors "github.com/chainedpixel/go-dte-signer/internal/domain/errors"
	"github.com/chainedpixel/go-dte-signer/internal/domain/models"
)

// KeyProcessor handles cryptographic key processing
type KeyProcessor struct{}

// NewKeyProcessor creates a new key processor
func NewKeyProcessor() *KeyProcessor {
	return &KeyProcessor{}
}

// BytesToPrivateKey converts a byte array to an RSA private key
func (k *KeyProcessor) BytesToPrivateKey(bytes []byte) (*models.Certificate, error) {
	priv, err := x509.ParsePKCS8PrivateKey(bytes)
	if err != nil {
		return nil, domainErrors.NewDomainError(err.Error(), domainErrors.CodeUncatalogued)
	}

	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, domainErrors.NewDomainError("key is not RSA type", domainErrors.CodeInvalid)
	}

	cert := &models.Certificate{
		DecodedPrivateKey: rsaPriv,
	}

	return cert, nil
}

// HashPassword hashes a password using SHA-512
func (k *KeyProcessor) HashPassword(password string) (string, error) {
	hasher := sha512.New()
	_, err := hasher.Write([]byte(password))
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	hashedBytes := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString, nil
}

// VerifyPassword verifies a password against a hashed password
func (k *KeyProcessor) VerifyPassword(plainPassword, hashedPassword string) (bool, error) {
	hashedInput, err := k.HashPassword(plainPassword)
	if err != nil {
		return false, err
	}
	return hashedInput == hashedPassword, nil
}
