package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"certforge-api/config"
	"certforge-api/models"
	"certforge-api/storage"
	"certforge-api/utils"
)

// CreateCA godoc
// @Summary      Create Certificate Authority
// @Description  Generate a new self-signed CA to sign certificates
// @Tags         Certificate Authority
// @Accept       json
// @Produce      json
// @Param        request body models.CARequest true "CA Configuration"
// @Success      200 {object} models.CAResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /ca/create [post]
func CreateCA(c *gin.Context) {
	var req models.CARequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Valori di default
	if req.ValidDays == 0 {
		req.ValidDays = 3650 // 10 anni
	}
	if req.KeySize == 0 {
		req.KeySize = 4096
	}

	// Genera CA
	cert, key, err := utils.GenerateCA(
		req.CommonName,
		req.Organization,
		req.Country,
		req.ValidDays,
		req.KeySize,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "generation_failed",
			Message: err.Error(),
		})
		return
	}

	// Salva CA globalmente
	config.SetCA(cert, key)

	// Salva su file
	certPEM := utils.EncodeCertToPEM(cert)
	keyPEM := utils.EncodeKeyToPEM(key)

	if err := storage.SaveCA(certPEM, keyPEM); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "save_failed",
			Message: err.Error(),
		})
		return
	}

	response := models.CAResponse{
		Message:      "CA creata con successo",
		CommonName:   cert.Subject.CommonName,
		SerialNumber: cert.SerialNumber.String(),
		NotBefore:    cert.NotBefore,
		NotAfter:     cert.NotAfter,
		Certificate:  string(certPEM),
		PrivateKey:   string(keyPEM),
	}

	c.JSON(http.StatusOK, response)
}

// GetCAInfo godoc
// @Summary Get CA information
// @Description Retrieve information about the current Certificate Authority
// @Tags Certificate Authority
// @Produce json
// @Success 200 {object} object{common_name=string,organization=[]string,country=[]string,serial_number=string,not_before=string,not_after=string,is_ca=bool,certificate=string}
// @Failure 404 {object} models.ErrorResponse
// @Router /ca/info [get]
func GetCAInfo(c *gin.Context) {
	cert, _ := config.GetCA()

	if cert == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "ca_not_found",
			Message: "CA non ancora creata",
		})
		return
	}

	certPEM := utils.EncodeCertToPEM(cert)

	response := gin.H{
		"common_name":   cert.Subject.CommonName,
		"organization":  cert.Subject.Organization,
		"country":       cert.Subject.Country,
		"serial_number": cert.SerialNumber.String(),
		"not_before":    cert.NotBefore,
		"not_after":     cert.NotAfter,
		"is_ca":         cert.IsCA,
		"certificate":   string(certPEM),
	}

	c.JSON(http.StatusOK, response)
}

// GetCACertificate godoc
// @Summary      Download CA certificate
// @Description  Download the CA certificate in PEM format (without private key)
// @Tags         Certificate Authority
// @Produce      application/x-pem-file
// @Success      200 {file} string "CA Certificate PEM"
// @Failure      404 {object} models.ErrorResponse
// @Router       /ca/certificate [get]
func GetCACertificate(c *gin.Context) {
	cert, _ := config.GetCA()

	if cert == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "ca_not_found",
			Message: "CA non ancora creata",
		})
		return
	}

	certPEM := utils.EncodeCertToPEM(cert)

	c.Header("Content-Type", "application/x-pem-file")
	c.Header("Content-Disposition", "attachment; filename=ca-cert.pem")
	c.Data(http.StatusOK, "application/x-pem-file", certPEM)
}
