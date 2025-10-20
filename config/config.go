package config

import (
	"crypto/rsa"
	"crypto/x509"
	"sync"
)

var (
	// CA globale
	CACert *x509.Certificate
	CAKey  *rsa.PrivateKey
	mu     sync.RWMutex
)

// SetCA imposta la CA globale
func SetCA(cert *x509.Certificate, key *rsa.PrivateKey) {
	mu.Lock()
	defer mu.Unlock()
	CACert = cert
	CAKey = key
}

// GetCA ottiene la CA globale
func GetCA() (*x509.Certificate, *rsa.PrivateKey) {
	mu.RLock()
	defer mu.RUnlock()
	return CACert, CAKey
}

// HasCA verifica se la CA esiste
func HasCA() bool {
	mu.RLock()
	defer mu.RUnlock()
	return CACert != nil && CAKey != nil
}
