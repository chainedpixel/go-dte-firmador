package ports

import (
	"context"

	"github.com/chainedpixel/go-dte-signer/internal/domain/models"
)

// SigningService defines the operations for signing documents
type SigningService interface {
	// SignDocument signs a document with the specified certificate
	SignDocument(ctx context.Context, request *models.CertificateRequest) (string, error)
}

// KeyProcessor defines operations for processing cryptographic keys
type KeyProcessor interface {
	// BytesToPrivateKey converts a byte array to an RSA private key
	BytesToPrivateKey(bytes []byte) (*models.Certificate, error)
}

// DocumentSigner defines operations for signing documents
type DocumentSigner interface {
	// Sign signs a document with the provided certificate
	Sign(ctx context.Context, certificate *models.Certificate, documentData interface{}) (string, error)
}
