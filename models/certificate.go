package models

import "time"

// CARequest - Richiesta per creare una CA
type CARequest struct {
	CommonName   string `json:"common_name" binding:"required"`
	Organization string `json:"organization" binding:"required"`
	Country      string `json:"country" binding:"required,len=2"`
	ValidDays    int    `json:"valid_days"`
	KeySize      int    `json:"key_size"`
}

// CAResponse - Risposta creazione CA
type CAResponse struct {
	Message      string    `json:"message"`
	CommonName   string    `json:"common_name"`
	SerialNumber string    `json:"serial_number"`
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	Certificate  string    `json:"certificate"`
	PrivateKey   string    `json:"private_key,omitempty"`
}

// CertRequest - Richiesta per creare un certificato
type CertRequest struct {
	CommonName   string   `json:"common_name" binding:"required"`
	Organization string   `json:"organization" binding:"required"`
	Country      string   `json:"country" binding:"required,len=2"`
	ValidDays    int      `json:"valid_days"`
	SANs         []string `json:"sans,omitempty"`
}

// CertResponse - Risposta creazione certificato
type CertResponse struct {
	ID            string    `json:"id"`
	CommonName    string    `json:"common_name"`
	Type          string    `json:"type"`
	Certificate   string    `json:"certificate"`
	PrivateKey    string    `json:"private_key"`
	CACertificate string    `json:"ca_certificate"`
	SerialNumber  string    `json:"serial_number"`
	NotBefore     time.Time `json:"not_before"`
	NotAfter      time.Time `json:"not_after"`
	DownloadURL   string    `json:"download_url"`
}

// CertInfo - Info certificato
type CertInfo struct {
	ID           string    `json:"id"`
	CommonName   string    `json:"common_name"`
	Type         string    `json:"type"`
	SerialNumber string    `json:"serial_number"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Revoked      bool      `json:"revoked"`
}

// ErrorResponse - Risposta errore
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
