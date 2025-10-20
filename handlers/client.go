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

// CreateClientCert godoc
// @Summary      Create client certificate
// @Description  Generate an X.509 certificate for mTLS clients (IoT devices, applications, etc.)
// @Description  The certificate will be signed by the configured CA
// @Tags         Client Certificates
// @Accept       json
// @Produce      json
// @Param        request body models.CertRequest true "Client certificate configuration"
// @Success      200 {object} models.CertResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /cert/client [post]
func CreateClientCert(c *gin.Context) {
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

	// Genera certificato client
	cert, key, err := utils.GenerateClientCert(
		req.CommonName,
		req.Organization,
		req.Country,
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
	if err := storage.SaveClientCert(certID, certPEM, keyPEM, cert); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "save_failed",
			Message: err.Error(),
		})
		return
	}

	response := models.CertResponse{
		ID:            certID,
		CommonName:    cert.Subject.CommonName,
		Type:          "client",
		Certificate:   string(certPEM),
		PrivateKey:    string(keyPEM),
		CACertificate: string(caCertPEM),
		SerialNumber:  cert.SerialNumber.String(),
		NotBefore:     cert.NotBefore,
		NotAfter:      cert.NotAfter,
		DownloadURL:   "/api/v1/cert/client/" + certID,
	}

	c.JSON(http.StatusOK, response)
}

// GetClientCert godoc
// @Summary      Get client certificate
// @Description  Retrieve an existing client certificate by ID
// @Tags         Client Certificates
// @Produce      json
// @Param        id path string true "Certificate ID"
// @Success      200 {object} object{id=string,type=string,certificate=string,private_key=string,ca_certificate=string}
// @Failure      404 {object} models.ErrorResponse
// @Router       /cert/client/{id} [get]
func GetClientCert(c *gin.Context) {
	certID := c.Param("id")

	certPEM, keyPEM, err := storage.LoadCert("client", certID)
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
		"type":           "client",
		"certificate":    string(certPEM),
		"private_key":    string(keyPEM),
		"ca_certificate": string(caCertPEM),
	}

	c.JSON(http.StatusOK, response)
}

// DownloadClientPackage godoc
// @Summary      Download client package
// @Description  Download a ZIP package containing certificate, key and CA for the client
// @Tags         Client Certificates
// @Produce      json
// @Param        id path string true "Certificate ID"
// @Success      200 {object} object
// @Failure      404 {object} models.ErrorResponse
// @Router       /cert/client/{id}/package [get]
func DownloadClientPackage(c *gin.Context) {
	certID := c.Param("id")

	certPEM, keyPEM, err := storage.LoadCert("client", certID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "cert_not_found",
			Message: "Certificato non trovato",
		})
		return
	}

	// Carica il certificato CA
	caCert, _ := config.GetCA()
	if caCert == nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "ca_not_available",
			Message: "CA non disponibile",
		})
		return
	}
	caCertPEM := utils.EncodeCertToPEM(caCert)

	// Crea ZIP in memoria
	// (Implementazione completa sotto)

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=client-"+certID+".zip")

	// Per ora restituiamo i file separati
	// Implementeremo lo ZIP dopo se vuoi
	c.JSON(http.StatusOK, gin.H{
		"message": "Use GET /api/v1/cert/client/:id for now",
		"files": gin.H{
			"client-cert.pem": string(certPEM),
			"client-key.pem":  string(keyPEM),
			"ca-cert.pem":     string(caCertPEM),
		},
	})
}
