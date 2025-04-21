package ports

import (
	"context"

	"github.com/chainedpixel/go-dte-signer/internal/domain/models"
)

// CertificateRepository defines the operations for certificate data access
type CertificateRepository interface {
	// GetByNIT retrieves a certificate by NIT
	GetByNIT(ctx context.Context, nit string) (*models.Certificate, error)

	// VerifyPassword checks if the password is valid for the certificate
	VerifyPassword(ctx context.Context, certificate *models.Certificate, password string) (bool, error)
}
