package models

import (
	"crypto/rsa"
	"encoding/base64"
)

// Certificate represents a digital certificate used for signing
type Certificate struct {
	ID                string          `json:"_id" xml:"_id"`
	Active            bool            `json:"activo" xml:"activo"`
	NIT               string          `json:"nit" xml:"nit"`
	PrivateKey        Key             `json:"privateKey" xml:"privateKey"`
	PublicKey         Key             `json:"publicKey" xml:"publicKey"`
	DecodedPrivateKey *rsa.PrivateKey `json:"-" xml:"-"` // Not mapped to XML/JSON, for internal use only
}

// Key represents a cryptographic key
type Key struct {
	Algorithm string `json:"algorithm" xml:"algorithm"`
	Password  string `json:"clave" xml:"clave"`
	Encoded   string `json:"encodied" xml:"encodied"`
	Format    string `json:"format" xml:"format"`
	KeyType   string `json:"keyType" xml:"keyType"`
}

// DecodePrivateKey decodes the base64-encoded private key
func (c *Certificate) DecodePrivateKey() ([]byte, error) {
	return base64.StdEncoding.DecodeString(c.PrivateKey.Encoded)
}

// IsActive returns whether the certificate is active
func (c *Certificate) IsActive() bool {
	return c.Active
}

// HasPrivateKey checks if the certificate has a private key
func (c *Certificate) HasPrivateKey() bool {
	return c.PrivateKey.Encoded != ""
}

// CertificateRequest contains the request data for signing a document
type CertificateRequest struct {
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

// Validate checks if the certificate request contains the required fields
func (r *CertificateRequest) Validate() bool {
	return r.NIT != "" && r.PrivateKeyPassword != "" && r.DocumentJSON != nil
}
