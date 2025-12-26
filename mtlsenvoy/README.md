## MTLS with envoy

1. Follow steps in certs/README.md to generate all the certs needed

## Curl
```
curl -vk https://localhost:8081/fast \
  --cert certs/client.crt \
  --key certs/client.key \
  --cacert certs/ca.crt

--cert certs/client.crt
Client certificate
Proves the identity of the client to the server. Envoy will check that this cert is signed by a trusted CA.
--key certs/client.key
Client private key
Corresponds to the client certificate. Used to sign the TLS handshake, proving ownership of the client cert.
--cacert certs/ca.crt
Certificate Authority
Tells curl which CA it should trust to verify the server’s certificate. Without this, curl would reject the server’s TLS certificate if it’s self-signed (which it is in your local setup).
```


### Curls that fail
```
curl -vk https://localhost:8081/fast --cacert certs/ca.crt
    - should fail

curl -vk https://localhost:8081/fast \
  --cert certs/bad-client.crt \
  --key certs/bad-client.key \
  --cacert certs/ca.crt

-> should fail with
error:1404C418:SSL routines:ST_OK:tlsv1 alert unknown ca, errno 0

curl https://localhost:8081/fast \
  --cert certs/client.crt \
  --key certs/client.key \
  --fail


curl https://localhost:8081/fast \
  --cert certs/client.crt \
  --key certs/client.key \
  --fail


should fail with
-> curl: (60) SSL certificate problem: unable to get local issuer certificate

curl -vk https://localhost:8081/fast \
  --cert certs/client.crt \
  --key certs/service-a.key \
  --cacert certs/ca.crt

should fail with
-> * unable to set private key file: 'certs/service-a.key' type PEM
```


## How envoy listeners + clusters work
```
Client
   |
   v
[Envoy Listener]  <-- accepts requests
   |
   |-- routing rules --> [Cluster A: service-a instances]
   |-- routing rules --> [Cluster B: service-b instances]
   v
Backends (actual services)
```

The listener is like the front door -> it then decides what to do and forwards it to the "cluster"


## Debugging
1. 	Know the paths to all relevant certificates and keys:
	•	client.crt / client.key
	•	service-a.crt / service-a.key
	•	service-b.crt / service-b.key
	•	ca.crt (your root CA)
	•	Know which Envoy listener and cluster configs correspond to each service.

2. 
- Verify that the client cert is valid and signed by the trusted CA:
```
openssl verify -CAfile certs/ca.crt certs/client.crt
```
- Inspect the cert details
```
openssl x509 -in certs/client.crt -text -noout
```
Check:
	•	Issuer matches CN=local-ca
	•	Validity dates are current
	•	Subject / SAN fields are as expected

3. Test the TLS handshake directly
```
openssl s_client -connect localhost:8081 -cert certs/client.crt -key certs/client.key -CAfile certs/ca.crt
```

4. Verify the server (listener) configuration
```
transport_socket:
  name: envoy.transport_sockets.tls
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext
    require_client_certificate: true
    common_tls_context:
      tls_certificates:
        - certificate_chain: { filename: "/etc/envoy/certs/service-a.crt" }
          private_key: { filename: "/etc/envoy/certs/service-a.key" }
      validation_context:
        trusted_ca: { filename: "/etc/envoy/certs/ca.crt" }
```

If envoy can't reach a backend
```
clusters:
  - name: service-b
    connect_timeout: 1s
    type: STRICT_DNS
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: service-b
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: service-b
                    port_value: 8080
    transport_socket:
      name: envoy.transport_sockets.tls
      typed_config:
        "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
        common_tls_context:
          tls_certificates:
            - certificate_chain: { filename: "/etc/envoy/certs/service-a.crt" }
              private_key: { filename: "/etc/envoy/certs/service-a.key" }
          validation_context:
            trusted_ca: { filename: "/etc/envoy/certs/ca.crt" }
```
```
 CLIENT (curl)
   -----------------
   | Client Cert    |
   | Client Key     |
   | CA Cert        |
   -----------------
           |
           | TLS handshake (mutual TLS)
           | Presents client cert to Envoy listener
           v
   ENVoy PROXY (Listener)
   -------------------------
   | Downstream TLS context |
   | require_client_cert: true
   | Server Cert / Key      |
   | Trusted CA for clients |
   -------------------------
           |
           | TLS handshake (Envoy verifies client)
           |
           | Route request to cluster (service-b)
           v
   BACKEND (service-b)
   -------------------------
   | Upstream TLS context   |
   | Envoy client cert      |
   | Backend Cert / Key     |
   | Trusted CA (for Envoy)|
   -------------------------
           ^
           | Envoy performs TLS handshake upstream
           |
```
