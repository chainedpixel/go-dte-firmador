package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chainedpixel/go-dte-signer/internal/domain/errors"
	"github.com/chainedpixel/go-dte-signer/internal/domain/models"
	"github.com/chainedpixel/go-dte-signer/internal/domain/ports"
)

// SigningService implements the ports.SigningService interface
type SigningService struct {
	certRepo       ports.CertificateRepository
	documentSigner ports.DocumentSigner
}

// NewSigningService creates a new signing service
func NewSigningService(certRepo ports.CertificateRepository, documentSigner ports.DocumentSigner) *SigningService {
	return &SigningService{
		certRepo:       certRepo,
		documentSigner: documentSigner,
	}
}

// SignDocument signs a document with the specified certificate
func (s *SigningService) SignDocument(ctx context.Context, request *models.CertificateRequest) (string, error) {
	// 1: Validate the request
	if !request.Validate() {
		return "", errors.NewRequiredDataError("required_data")
	}

	// 2: Retrieve the certificate by NIT
	certificate, err := s.certRepo.GetByNIT(ctx, request.NIT)
	if err != nil {
		return "", err
	}

	// 3: Verify the private key password
	valid, err := s.certRepo.VerifyPassword(ctx, certificate, request.PrivateKeyPassword)
	if err != nil {
		return "", err
	}

	// 4: Check if the password is valid
	if !valid {
		return "", errors.NewPasswordInvalidError(fmt.Sprintf("password_invalid", request.NIT))
	}

	// 5: Process the document JSON
	var documentData []byte
	switch v := request.DocumentJSON.(type) {
	case string:
		documentData = []byte(v)
	case []byte:
		documentData = v
	default:
		// Marshal the document JSON to a string
		var err error
		documentData, err = json.Marshal(request.DocumentJSON)
		if err != nil {
			return "", errors.NewDomainError("json_to_string_conversion", errors.CodeJSONToStrConversion)
		}
	}

	// 6: Sign the document
	signedJWS, err := s.documentSigner.Sign(ctx, certificate, documentData)
	if err != nil {
		return "", err
	}

	// 7: Return the signed JWS
	return signedJWS, nil
}
