package cypher

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"github.com/go-jose/go-jose/v3"

	domainErrors "github.com/chainedpixel/go-dte-signer/internal/domain/errors"
	"github.com/chainedpixel/go-dte-signer/internal/domain/models"
)

// JWSSigner handles JWS signing operations
type JWSSigner struct{}

// NewJWSSigner creates a new JWS signer
func NewJWSSigner() *JWSSigner {
	return &JWSSigner{}
}

// Sign signs a document with the provided certificate
func (s *JWSSigner) Sign(ctx context.Context, certificate *models.Certificate, documentData interface{}) (string, error) {
	// Ensure the private key is available
	if certificate.DecodedPrivateKey == nil {
		return "", domainErrors.NewDomainError("private key not available", domainErrors.CodeInvalid)
	}

	// Convert document data to a string if necessary
	var documentStr string
	switch v := documentData.(type) {
	case string:
		documentStr = v
	case []byte:
		documentStr = string(v)
	default:
		jsonData, err := json.Marshal(documentData)
		if err != nil {
			return "", domainErrors.NewDomainError(err.Error(), domainErrors.CodeJSONToStrConversion)
		}
		documentStr = string(jsonData)
	}

	// Create signer with RSA-SHA512 algorithm
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.RS512,
		Key:       certificate.DecodedPrivateKey,
	}, nil)
	if err != nil {
		return "", domainErrors.NewDomainError(err.Error(), domainErrors.CodeUncatalogued)
	}

	// Sign the document
	object, err := signer.Sign([]byte(documentStr))
	if err != nil {
		return "", domainErrors.NewDomainError(err.Error(), domainErrors.CodeUncatalogued)
	}

	// Serialize the signature
	serialized, err := object.CompactSerialize()
	if err != nil {
		return "", domainErrors.NewDomainError(err.Error(), domainErrors.CodeUncatalogued)
	}

	return serialized, nil
}

// SignWithPrivateKey signs data with a raw private key
func (s *JWSSigner) SignWithPrivateKey(data string, privateKey *rsa.PrivateKey) (string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.RS512,
		Key:       privateKey,
	}, nil)
	if err != nil {
		return "", err
	}

	object, err := signer.Sign([]byte(data))
	if err != nil {
		return "", err
	}

	serialized, err := object.CompactSerialize()
	if err != nil {
		return "", err
	}

	return serialized, nil
}
