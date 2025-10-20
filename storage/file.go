package storage

import (
	// "crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	certsDir  = "./certs"
	caDir     = certsDir + "/ca"
	serverDir = certsDir + "/servers"
	clientDir = certsDir + "/clients"
	metaFile  = "metadata.json"
)

// CertMetadata - Metadati del certificato
type CertMetadata struct {
	ID           string    `json:"id"`
	CommonName   string    `json:"common_name"`
	Type         string    `json:"type"`
	SerialNumber string    `json:"serial_number"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Revoked      bool      `json:"revoked"`
}

func init() {
	os.MkdirAll(caDir, 0700)
	os.MkdirAll(serverDir, 0700)
	os.MkdirAll(clientDir, 0700)
}

// SaveCA salva la CA
func SaveCA(certPEM, keyPEM []byte) error {
	if err := os.WriteFile(filepath.Join(caDir, "ca-cert.pem"), certPEM, 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(caDir, "ca-key.pem"), keyPEM, 0600)
}

// LoadCA carica la CA
func LoadCA() ([]byte, []byte, error) {
	certPEM, err := os.ReadFile(filepath.Join(caDir, "ca-cert.pem"))
	if err != nil {
		return nil, nil, err
	}
	keyPEM, err := os.ReadFile(filepath.Join(caDir, "ca-key.pem"))
	if err != nil {
		return nil, nil, err
	}
	return certPEM, keyPEM, nil
}

// SaveServerCert salva un certificato server
func SaveServerCert(id string, certPEM, keyPEM []byte, cert *x509.Certificate) error {
	dir := filepath.Join(serverDir, id)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "cert.pem"), certPEM, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "key.pem"), keyPEM, 0600); err != nil {
		return err
	}

	// Salva metadata
	meta := CertMetadata{
		ID:           id,
		CommonName:   cert.Subject.CommonName,
		Type:         "server",
		SerialNumber: cert.SerialNumber.String(),
		IssuedAt:     cert.NotBefore,
		ExpiresAt:    cert.NotAfter,
		Revoked:      false,
	}
	return saveMetadata(dir, meta)
}

// SaveClientCert salva un certificato client
func SaveClientCert(id string, certPEM, keyPEM []byte, cert *x509.Certificate) error {
	dir := filepath.Join(clientDir, id)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "cert.pem"), certPEM, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "key.pem"), keyPEM, 0600); err != nil {
		return err
	}

	// Salva metadata
	meta := CertMetadata{
		ID:           id,
		CommonName:   cert.Subject.CommonName,
		Type:         "client",
		SerialNumber: cert.SerialNumber.String(),
		IssuedAt:     cert.NotBefore,
		ExpiresAt:    cert.NotAfter,
		Revoked:      false,
	}
	return saveMetadata(dir, meta)
}

// LoadCert carica un certificato
func LoadCert(certType, id string) ([]byte, []byte, error) {
	var dir string
	if certType == "server" {
		dir = filepath.Join(serverDir, id)
	} else {
		dir = filepath.Join(clientDir, id)
	}

	certPEM, err := os.ReadFile(filepath.Join(dir, "cert.pem"))
	if err != nil {
		return nil, nil, err
	}
	keyPEM, err := os.ReadFile(filepath.Join(dir, "key.pem"))
	if err != nil {
		return nil, nil, err
	}
	return certPEM, keyPEM, nil
}

// ListCerts lista tutti i certificati
func ListCerts() ([]CertMetadata, error) {
	var certs []CertMetadata

	// Lista server certs
	serverCerts, err := listCertsInDir(serverDir)
	if err == nil {
		certs = append(certs, serverCerts...)
	}

	// Lista client certs
	clientCerts, err := listCertsInDir(clientDir)
	if err == nil {
		certs = append(certs, clientCerts...)
	}

	return certs, nil
}

func listCertsInDir(dir string) ([]CertMetadata, error) {
	var certs []CertMetadata

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		metaPath := filepath.Join(dir, entry.Name(), metaFile)
		data, err := os.ReadFile(metaPath)
		if err != nil {
			continue
		}

		var meta CertMetadata
		if err := json.Unmarshal(data, &meta); err != nil {
			continue
		}

		certs = append(certs, meta)
	}

	return certs, nil
}

func saveMetadata(dir string, meta CertMetadata) error {
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, metaFile), data, 0644)
}
