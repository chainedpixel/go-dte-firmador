package adapters

import (
	"context"
	"encoding/xml"
	"os"
	"path/filepath"

	domainErrors "github.com/chainedpixel/go-dte-signer/internal/domain/errors"
	"github.com/chainedpixel/go-dte-signer/internal/domain/models"
	"github.com/chainedpixel/go-dte-signer/internal/infrastructure/cypher"
	"github.com/chainedpixel/go-dte-signer/pkg/logs"
)

// FileCertificateRepository implements a file-based certificate repository
type FileCertificateRepository struct {
	basePath     string
	keyProcessor *cypher.KeyProcessor
}

// NewFileCertificateRepository creates a new file-based certificate repository
func NewFileCertificateRepository(basePath string, keyProcessor *cypher.KeyProcessor) *FileCertificateRepository {
	return &FileCertificateRepository{
		basePath:     basePath,
		keyProcessor: keyProcessor,
	}
}

// GetByNIT retrieves a certificate by NIT
func (r *FileCertificateRepository) GetByNIT(ctx context.Context, nit string) (*models.Certificate, error) {
	// Construct the file path
	filePath := filepath.Join(r.basePath, nit+".crt")

	// Read the certificate file
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logs.Error("Certificate does not exist")
			return nil, domainErrors.NewDomainError(filePath, domainErrors.CodeFileNotFound)
		}
		logs.Error("Failed to read certificate:", err)
		return nil, domainErrors.NewDomainError(err.Error(), domainErrors.CodeUncatalogued)
	}

	// Parse the XML content
	var certificate models.Certificate
	if err = xml.Unmarshal(content, &certificate); err != nil {
		logs.Error("Failed to unmarshal XML:", err)
		return nil, domainErrors.NewDomainError(err.Error(), domainErrors.CodeUncatalogued)
	}

	// Check if the certificate is active
	if !certificate.IsActive() {
		logs.Error("Certificate is not active")
		return nil, domainErrors.NewDomainError("cert_not_found", domainErrors.CodeCertNotFound)
	}

	// Check if the certificate has a private key
	if !certificate.HasPrivateKey() {
		logs.Error("Certificate does not have a private key")
		return nil, domainErrors.NewDomainError("invalid", domainErrors.CodeInvalid)
	}

	// Decode the private key
	decodedBytes, err := certificate.DecodePrivateKey()
	if err != nil {
		logs.Error("Failed to decode private key:", err)
		return nil, domainErrors.NewDomainError(err.Error(), domainErrors.CodeJSONToStrConversion)
	}

	// Parse the private key
	certWithKey, err := r.keyProcessor.BytesToPrivateKey(decodedBytes)
	if err != nil {
		logs.Error("Failed to parse private key:", err)
		return nil, domainErrors.NewDomainError(err.Error(), domainErrors.CodeNoPublicKey)
	}

	// Update the certificate with the decoded private key
	certificate.DecodedPrivateKey = certWithKey.DecodedPrivateKey

	return &certificate, nil
}

// VerifyPassword checks if the password is valid for the certificate
func (r *FileCertificateRepository) VerifyPassword(ctx context.Context, certificate *models.Certificate, password string) (bool, error) {
	// Verify the password
	valid, err := r.keyProcessor.VerifyPassword(password, certificate.PrivateKey.Password)
	if err != nil {
		logs.Error("Failed to verify password:", err)
		return false, domainErrors.NewDomainError(err.Error(), domainErrors.CodeUncatalogued)
	}

	return valid, nil
}
