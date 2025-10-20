package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"certforge-api/config"
	"certforge-api/models"
	"certforge-api/storage"
	"certforge-api/utils"
)

// CreateServerCert godoc
// @Summary      Create server certificate
// @Description  Generate an X.509 certificate for TLS servers (web servers, MQTT brokers, etc.)
// @Description  The certificate will be signed by the configured CA
// @Tags         Server Certificates
// @Accept       json
// @Produce      json
// @Param        request body models.CertRequest true "Server certificate configuration"
// @Success      200 {object} models.CertResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /cert/server [post]
func CreateServerCert(c *gin.Context) {
	var req models.CertRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Verifica che la CA esista
	caCert, caKey := config.GetCA()
	if caCert == nil || caKey == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "ca_not_available",
			Message: "CA non disponibile. Creare prima la CA con POST /api/v1/ca/create",
		})
		return
	}

	// Valori di default
	if req.ValidDays == 0 {
		req.ValidDays = 365
	}

	// Genera certificato server
	cert, key, err := utils.GenerateServerCert(
		req.CommonName,
		req.Organization,
		req.Country,
		req.SANs,
		req.ValidDays,
		caCert,
		caKey,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "generation_failed",
			Message: err.Error(),
		})
		return
	}

	// Converti in PEM
	certPEM := utils.EncodeCertToPEM(cert)
	keyPEM := utils.EncodeKeyToPEM(key)
	caCertPEM := utils.EncodeCertToPEM(caCert)

	// Genera ID univoco
	certID := uuid.New().String()

	// Salva su storage
	if err := storage.SaveServerCert(certID, certPEM, keyPEM, cert); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "save_failed",
			Message: err.Error(),
		})
		return
	}

	response := models.CertResponse{
		ID:            certID,
		CommonName:    cert.Subject.CommonName,
		Type:          "server",
		Certificate:   string(certPEM),
		PrivateKey:    string(keyPEM),
		CACertificate: string(caCertPEM),
		SerialNumber:  cert.SerialNumber.String(),
		NotBefore:     cert.NotBefore,
		NotAfter:      cert.NotAfter,
		DownloadURL:   "/api/v1/cert/server/" + certID,
	}

	c.JSON(http.StatusOK, response)
}

// GetServerCert godoc
// @Summary      Get server certificate
// @Description  Retrieve an existing server certificate by ID
// @Tags         Server Certificates
// @Produce      json
// @Param        id path string true "Certificate ID"
// @Success      200 {object} object{id=string,type=string,certificate=string,private_key=string,ca_certificate=string}
// @Failure      404 {object} models.ErrorResponse
// @Router       /cert/server/{id} [get]
func GetServerCert(c *gin.Context) {
	certID := c.Param("id")

	certPEM, keyPEM, err := storage.LoadCert("server", certID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "cert_not_found",
			Message: "Certificato non trovato",
		})
		return
	}

	// Carica anche il certificato CA
	caCert, _ := config.GetCA()
	var caCertPEM []byte
	if caCert != nil {
		caCertPEM = utils.EncodeCertToPEM(caCert)
	}

	response := gin.H{
		"id":             certID,
		"type":           "server",
		"certificate":    string(certPEM),
		"private_key":    string(keyPEM),
		"ca_certificate": string(caCertPEM),
	}

	c.JSON(http.StatusOK, response)
}
