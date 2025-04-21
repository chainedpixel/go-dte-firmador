package usecases

import (
	"context"
	"errors"
	"fmt"

	errPackage "github.com/chainedpixel/go-dte-signer/internal/domain/errors"
	"github.com/chainedpixel/go-dte-signer/internal/domain/models"
	"github.com/chainedpixel/go-dte-signer/internal/domain/ports"
	"github.com/chainedpixel/go-dte-signer/pkg/i18n"
	"github.com/chainedpixel/go-dte-signer/pkg/response"
)

// DocumentSigningUseCase handles document signing operations
type DocumentSigningUseCase struct {
	signingService ports.SigningService
	translator     *i18n.Translator
}

// NewDocumentSigningUseCase creates a new document signing use case
func NewDocumentSigningUseCase(signingService ports.SigningService, translator *i18n.Translator) *DocumentSigningUseCase {
	return &DocumentSigningUseCase{
		signingService: signingService,
		translator:     translator,
	}
}

// DocumentSigningInput represents the input for document signing
type DocumentSigningInput struct {
	PublicKeyPassword  string      `json:"passwordPub"`
	PrivateKeyPassword string      `json:"passwordPri"`
	NIT                string      `json:"nit"`
	DocumentName       string      `json:"nombreDocumento"`
	SignatureName      string      `json:"nombreFirma"`
	CompactSerialized  string      `json:"compactSerialization"`
	DocumentJSON       interface{} `json:"dteJson"`
	Document           string      `json:"dte"`
	Active             bool        `json:"activo"`
	Path               string      `json:"path"`
}

// Execute processes a document signing request
func (uc *DocumentSigningUseCase) Execute(ctx context.Context, input DocumentSigningInput) (*response.Response, error) {
	// 1. Validate input
	if !uc.validateInput(input) {
		err := errPackage.NewRequiredDataError(uc.translator.T("required_data"))
		return uc.createErrorResponse(err), nil
	}

	// 2. Map input to domain model
	request := &models.CertificateRequest{
		PublicKeyPassword:  input.PublicKeyPassword,
		PrivateKeyPassword: input.PrivateKeyPassword,
		NIT:                input.NIT,
		DocumentName:       input.DocumentName,
		SignatureName:      input.SignatureName,
		CompactSerialized:  input.CompactSerialized,
		DocumentJSON:       input.DocumentJSON,
		Document:           input.Document,
		Active:             input.Active,
		Path:               input.Path,
	}

	// 3. Call domain service to sign document
	jws, err := uc.signingService.SignDocument(ctx, request)
	if err != nil {
		return uc.createErrorResponse(err), nil
	}

	return response.NewSuccessResponse(jws), nil
}

// validateInput validates the document signing input
func (uc *DocumentSigningUseCase) validateInput(input DocumentSigningInput) bool {
	return input.NIT != "" && input.PrivateKeyPassword != "" && input.DocumentJSON != nil
}

// createErrorResponse creates an error response from a domain error
func (uc *DocumentSigningUseCase) createErrorResponse(err error) *response.Response {
	var domainErr errPackage.DomainError
	ok := errors.As(err, &domainErr)
	if !ok {
		// Default to internal server error
		return response.NewErrorResponse("500", uc.translator.T("internal_server_error"))
	}

	// Translate the error message
	translatedMsg := uc.translator.T(domainErr.Message)
	if domainErr.Code == errPackage.CodePasswordInvalid {
		// Special case for password invalid, which needs formatting
		translatedMsg = fmt.Sprintf(translatedMsg, domainErr.Message)
	}

	return response.NewErrorResponse(domainErr.Code, translatedMsg)
}
