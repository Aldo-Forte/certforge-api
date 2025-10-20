package handlers

import (
	"certforge-api/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListCertificates godoc
// @Summary      List certificates
// @Description  Get a list of all generated certificates (server and client)
// @Tags         Certificates
// @Produce      json
// @Success      200 {object} object{total=int,certificates=[]storage.CertMetadata}
// @Failure      500 {object} object{error=string}
// @Router       /certs [get]
func ListCertificates(c *gin.Context) {
	certs, err := storage.ListCerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Errore nel recupero certificati",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":        len(certs),
		"certificates": certs,
	})
}

// RevokeCertificate godoc
// @Summary      Revoke certificate
// @Description  Revoke a certificate (not fully implemented yet)
// @Description  In production, this should update metadata, generate CRL and implement OCSP
// @Tags         Certificates
// @Produce      json
// @Param        id path string true "Certificate ID"
// @Success      501 {object} object{message=string,cert_id=string,todo=[]string}
// @Router       /cert/revoke/{id} [delete]
func RevokeCertificate(c *gin.Context) {
	certID := c.Param("id")

	// Simple implementation for now
	// In production you should:
	// - Update metadata.json with revoked: true
	// - Generate CRL (Certificate Revocation List)
	// - Implement OCSP (Online Certificate Status Protocol)

	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Revoca certificati non ancora implementata",
		"cert_id": certID,
		"todo": []string{
			"Implementare CRL",
			"Implementare OCSP",
			"Aggiornare database certificati",
		},
	})
}
