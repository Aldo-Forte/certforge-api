## CERTFORGE API
REST API per la generazione e gestione di certificati SSL/TLS autofirmati per ambienti di sviluppo, QA e testing.  
Ideale per configurare infrastrutture IoT di test, servizi MQTT, microservizi in staging e qualsiasi scenario che richieda comunicazioni sicure 
con autenticazione reciproca senza la necessità di una CA pubblica.
  
__Nota Importante__: 
Questa API genera certificati autofirmati destinati esclusivamente ad ambienti di sviluppo, test e QA. 
NON utilizzare in produzione. Per ambienti produttivi, utilizzare certificati rilasciati da 
Certificate Authority riconosciute (Let's Encrypt, DigiCert, ecc.).


## Caratteristiche
Certificate Authority (CA) Autofirmata: Crea e gestisci la tua CA privata per ambienti di test
Certificati Server: Genera certificati per server con supporto SAN (Subject Alternative Names)
Certificati Client: Crea certificati per dispositivi IoT, applicazioni e client in ambiente di test
mTLS Ready: Perfetto per implementare e testare autenticazione reciproca (mutual TLS)
Storage Persistente: I certificati vengono salvati automaticamente sul filesystem
Pacchetti Client: Download completo di certificato, chiave privata e CA in un unico pacchetto
API RESTful: Interfaccia semplice e intuitiva
Documentazione Swagger: API docs interattiva integrata
Format Flessibili: Supporto per PEM, PKCS12 e altri formati

## Quick Start
### Prerequisiti
Go 1.25 o superiore

## Documentazione
Accedi alla documentazione interattiva Swagger su: http://localhost:8080/swagger/index.html

## Utilizzo
### Crea una Certificate Authority (autofirmata)

```
curl -X POST http://localhost:8080/api/v1/ca/create \
  -H "Content-Type: application/json" \
  -d '{
    "common_name": "My Test CA",
    "organization": "My Company QA",
    "country": "IT",
    "valid_days": 3650,
    "key_size": 4096
  }'
```

## Genera un certificato server

```
curl -X POST http://localhost:8080/api/v1/cert/server \
  -H "Content-Type: application/json" \
  -d '{
    "common_name": "mqtt-test.example.com",
    "organization": "My Company",
    "country": "IT",
    "valid_days": 365,
    "sans": ["mqtt-test.example.com", "*.test.example.com"]
  }'
```


