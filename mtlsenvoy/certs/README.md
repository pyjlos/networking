## Certs
We need 3 types of certificates
- Certificate Authority
    - Signs all the certificates
    - this is our "trust" anchor
    - each service will only trust certs signed by this CA
- service A cert -> signed by CA
- service B cert -> signed by CA

### Creating the CA
```
# 1. Generate private key for CA
openssl genrsa -out ca.key 2048

# 2. Generate self-signed CA certificate
openssl req -x509 -new -nodes -key ca.key -sha256 -days 365 \
  -out ca.crt -subj "/CN=local-ca"
```

### Creating Service A Cert
- also do the same for service B
- you also need one for the client (curl from your local machine in this case)
```
# 1. Generate private key
openssl genrsa -out service-a.key 2048

# 2. Generate CSR (Certificate Signing Request)
openssl req -new -key service-a.key -out service-a.csr -subj "/CN=service-a"

# 3. Sign certificate with CA
openssl x509 -req -in service-a.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out service-a.crt -days 365 -sha256
```


### Validating the Certs
- openssl x509 -in service-a.crt -text -noout
    - Should see something like
```
Issuer: CN=local-ca
Subject: CN=service-a
```

### Creating a bad cert
```
openssl req -x509 -newkey rsa:2048 \
  -keyout certs/bad-client.key \
  -out certs/bad-client.crt \
  -days 365 -nodes \
  -subj "/CN=bad-client"

curl -vk https://localhost:8081/fast \
  --cert certs/bad-client.crt \
  --key certs/bad-client.key \
  --cacert certs/ca.crt
```