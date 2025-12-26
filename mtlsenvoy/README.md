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