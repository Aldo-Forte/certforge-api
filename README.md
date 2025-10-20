## CERTFORGE API (Versione Beta)
REST API per la generazione e gestione di certificati SSL/TLS autofirmati per ambienti di sviluppo, QA, testing e scopi didattici. Ideale per configurare infrastrutture IoT di test, servizi MQTT, microservizi in staging e per apprendere i concetti di PKI,  
crittografia e autenticazione mTLS senza la necessità di una CA pubblica.
  
__Nota Importante__: 
Questa API genera certificati autofirmati destinati esclusivamente ad ambienti di sviluppo, test e QA. 
NON utilizzare in produzione. Per ambienti produttivi, utilizzare certificati rilasciati da 
Certificate Authority riconosciute (Let's Encrypt, DigiCert, ecc.).


## Caratteristiche
__Certificate Authority (CA) Autofirmata__: Crea e gestisci la tua CA privata per ambienti di test e apprendimento  
__Certificati Server__: Genera certificati per server con supporto SAN (Subject Alternative Names)  
__Certificati Client__: Crea certificati per dispositivi IoT, applicazioni e client in ambiente di test  
__mTLS Ready__: Perfetto per implementare e testare autenticazione reciproca (mutual TLS)  
__Storage Persistente__: I certificati vengono salvati automaticamente sul filesystem  
__Pacchetti Client__: Download completo di certificato, chiave privata e CA in un unico pacchetto  
__API RESTful__: Interfaccia semplice e intuitiva  
__Documentazione Swagger__: API docs interattiva integrata  
__Format Flessibili__: Supporto per PEM, PKCS12 e altri formati  
__Didattico__: Perfetto per comprendere il funzionamento di PKI e certificati digitali   

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

## Genera un certificato client
```
curl -X POST http://localhost:8080/api/v1/cert/client \
  -H "Content-Type: application/json" \
  -d '{
    "common_name": "device-001",
    "organization": "My Company",
    "country": "IT",
    "valid_days": 365
  }'
```

##  Scarica il pacchetto completo del client
```
curl -X GET http://localhost:8080/api/v1/cert/client/{cert-id}/package \
  --output client-package.zip
```

## Il pacchetto ZIP contiene:
client-cert.pem - Certificato client  
client-key.pem - Chiave privata client  
ca-cert.pem - Certificato CA per la verifica  
client.p12 - Bundle PKCS12 (opzionale)  


## Struttura del progetto
```
certforge-api/
├── api/           # Definizioni API
├── certs/         # Storage certificati CA
├── config/        # Configurazione applicazione
├── docs/          # Documentazione Swagger generata
├── handlers/      # Handler HTTP
├── middleware/    # Middleware Gin
├── models/        # Modelli dati
├── output/        # Certificati generati
├── storage/       # Layer di persistenza
├── utils/         # Utility e helper
└── main.go        # Entry point
```


## Casi d'uso
Testing e Sviluppo
IoT Device Testing: Genera certificati unici per ogni dispositivo IoT in ambiente di test  
MQTT Broker (Staging/QA): Configura broker MQTT con autenticazione mTLS per testing  
Microservizi (Dev/Staging): Test di comunicazione sicura tra microservizi   
API Gateway Testing: Test di autenticazione client basata su certificati  
Ambiente di Sviluppo Locale: SSL/TLS per sviluppo in locale senza warning del browser  
CI/CD Pipeline: Certificati per test automatizzati di connessioni sicure  
Proof of Concept: Validazione veloce di architetture mTLS  

## Scopi Didattici e Formativi  
Corsi di Sicurezza Informatica: Strumento pratico per insegnare PKI e crittografia  
Workshop mTLS: Dimostrazioni pratiche di autenticazione reciproca  
Laboratori Universitari: Esercitazioni su certificati digitali e X.509  
Formazione Aziendale: Training su SSL/TLS per team di sviluppo  
Tutorial e Blog: Esempi concreti per articoli tecnici  
Apprendimento Self-Service: Sperimentare con certificati in ambiente sicuro  

## Valore Didattico
Questo progetto è particolarmente utile per:  
Comprendere la gerarchia dei certificati (CA → Server/Client)  
Visualizzare la struttura X.509 e i suoi campi  
Sperimentare con Subject Alternative Names (SANs)  
Implementare e testare mTLS (mutual TLS authentication)  
Capire il funzionamento di chiavi pubbliche/private  
Esplorare i formati dei certificati (PEM, DER, PKCS12)  
Simulare scenari reali di PKI enterprise  

## Sicurezza
ATTENZIONE! Solo per ambienti non produttivi: I certificati sono autofirmati e non riconosciuti da CA pubbliche  
Le chiavi private della CA sono salvate su filesystem locale  
Supporto per chiavi RSA da 2048 a 4096 bit  
Certificati conformi agli standard X.509  
Best practice per la generazione di certificati di test  

## Configurazione
I certificati vengono salvati nelle seguenti directory:  
./certs/ - Certificato e chiave CA  
./output/server/ - Certificati server  
./output/client/ - Certificati client  
 
## Contribuire
I contributi sono benvenuti!  

Fai un fork del progetto  
Crea un branch per la tua feature (git checkout -b feature/AmazingFeature)  
Committa le modifiche (git commit -m 'Add some AmazingFeature')  
Pusha sul branch (git push origin feature/AmazingFeature)  
Apri una Pull Request  

## Licenza  
Questo progetto è distribuito sotto licenza MIT. Vedi il file LICENSE per maggiori dettagli.  
Autore: Aldo Forte  
Email: software@aldoforte.it  
  
## Ringraziamenti  
[Gin Web Framework](https://gin-gonic.com)  
Swaggo  
Go Crypto Package  
 
Se questo progetto ti è stato utile per testing o apprendimento, considera di lasciare una stella!  





