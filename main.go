package main

import (
	"log"

	"certforge-api/config"
	_ "certforge-api/docs"
	"certforge-api/handlers"
	"certforge-api/storage"
	"certforge-api/utils"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @title Certificate API
// @version 1.0
// @description REST API for generating SSL/TLS and mTLS certificates
// @description Generate Certificate Authority, server and client certificates for TLS/mTLS configurations

// @contact.name   Aldo Forte
// @contact.email  software@aldoforte.it

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes http https
func main() {
	// Prova a caricare la CA esistente
	loadExistingCA()

	// Crea router Gin
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Gruppo API v1
	v1 := r.Group("/api/v1")
	{
		// Endpoint CA
		ca := v1.Group("/ca")
		{
			ca.POST("/create", handlers.CreateCA)
			ca.GET("/info", handlers.GetCAInfo)
			ca.GET("/certificate", handlers.GetCACertificate)
		}

		// Endpoint certificati
		cert := v1.Group("/cert")
		{
			// Server certificates
			cert.POST("/server", handlers.CreateServerCert)
			cert.GET("/server/:id", handlers.GetServerCert)

			// Client certificates
			cert.POST("/client", handlers.CreateClientCert)
			cert.GET("/client/:id", handlers.GetClientCert)
			cert.GET("/client/:id/package", handlers.DownloadClientPackage)

			// Revoca (non implementata completamente)
			cert.DELETE("/revoke/:id", handlers.RevokeCertificate)
		}

		// Lista certificati
		v1.GET("/certs", handlers.ListCertificates)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"ca_loaded": config.HasCA(),
		})
	})

	// Root endpoint con documentazione
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "Certificate API",
			"version":     "1.0.0",
			"description": "API REST per generazione certificati SSL/TLS e mTLS",
			"endpoints": gin.H{
				"health": "GET /health",
				"ca": gin.H{
					"create":      "POST /api/v1/ca/create",
					"info":        "GET /api/v1/ca/info",
					"certificate": "GET /api/v1/ca/certificate",
				},
				"certificates": gin.H{
					"create_server": "POST /api/v1/cert/server",
					"create_client": "POST /api/v1/cert/client",
					"get_server":    "GET /api/v1/cert/server/:id",
					"get_client":    "GET /api/v1/cert/client/:id",
					"download":      "GET /api/v1/cert/client/:id/package",
					"list":          "GET /api/v1/certs",
					"revoke":        "DELETE /api/v1/cert/revoke/:id",
				},
			},
			"examples": gin.H{
				"create_ca": gin.H{
					"method": "POST",
					"url":    "/api/v1/ca/create",
					"body": gin.H{
						"common_name":  "Mia CA",
						"organization": "La Mia Azienda",
						"country":      "IT",
						"valid_days":   3650,
						"key_size":     4096,
					},
				},
				"create_server": gin.H{
					"method": "POST",
					"url":    "/api/v1/cert/server",
					"body": gin.H{
						"common_name":  "mqtt.example.com",
						"organization": "Example Inc",
						"country":      "IT",
						"valid_days":   365,
						"sans":         []string{"mqtt.example.com", "*.example.com"},
					},
				},
				"create_client": gin.H{
					"method": "POST",
					"url":    "/api/v1/cert/client",
					"body": gin.H{
						"common_name":  "device-001",
						"organization": "Example Inc",
						"country":      "IT",
						"valid_days":   365,
					},
				},
			},
		})
	})

	log.Println("🚀 Certificate API Server")
	log.Println("📡 Listening on http://localhost:8080")
	log.Println("📚 Documentation: http://localhost:8080")
	log.Println("")

	if config.HasCA() {
		log.Println("✅ CA loaded successfully")
	} else {
		log.Println("⚠️  No CA found. Create one with POST /api/v1/ca/create")
	}
	log.Println("")

	r.Run(":8080")
}

// loadExistingCA prova a caricare una CA esistente
func loadExistingCA() {
	certPEM, keyPEM, err := storage.LoadCA()
	if err != nil {
		return
	}

	cert, err := utils.DecodeCertFromPEM(certPEM)
	if err != nil {
		log.Printf("⚠️  Error decoding CA certificate: %v\n", err)
		return
	}

	key, err := utils.DecodeKeyFromPEM(keyPEM)
	if err != nil {
		log.Printf("⚠️  Error decoding CA key: %v\n", err)
		return
	}

	config.SetCA(cert, key)
}
