## Reverse Proxy using nginx as lb

### Getting started
```
docker compose up -d
docker ps -> should see 3 backend containers and one nginx container
docker logs -f nginx
```

### Sending requests
```
curl localhost:8080/fast
curl localhost:8080/fast
curl localhost:8080/fast
curl localhost:8080/fast
-> should see these requests round-robin
192.168.65.1 - GET /fast HTTP/1.1 status=200 backend=backend-1 upstream=172.18.0.3:8080 request_time=0.001 upstream_time=0.001
192.168.65.1 - GET /fast HTTP/1.1 status=200 backend=backend-2 upstream=172.18.0.2:8080 request_time=0.001 upstream_time=0.001
192.168.65.1 - GET /fast HTTP/1.1 status=200 backend=backend-3 upstream=172.18.0.4:8080 request_time=0.001 upstream_time=0.001
192.168.65.1 - GET /fast HTTP/1.1 status=200 backend=backend-1 upstream=172.18.0.3:8080 request_time=0.001 upstream_time=0.001
192.168.65.1 - GET /fast HTTP/1.1 status=200 backend=backend-2 upstream=172.18.0.2:8080 request_time=0.001 upstream_time=0.001
192.168.65.1 - GET /fast HTTP/1.1 status=200 backend=backend-3 upstream=172.18.0.4:8080 request_time=0.001 upstream_time=0.001
192.168.65.1 - GET /fast HTTP/1.1 status=200 backend=backend-1 upstream=172.18.0.3:8080 request_time=0.001 upstream_time=0.001
192.168.65.1 - GET /fast HTTP/1.1 status=200 backend=backend-2 upstream=172.18.0.2:8080 request_time=0.000 upstream_time=0.001
192.168.65.1 - GET /fast HTTP/1.1 status=200 backend=backend-3 upstream=172.18.0.4:8080 request_time=0.001 upstream_time=0.001
```

### Mimic timeouts
```
curl localhost:8080/slow --max-time 2
curl localhost:8080/slow --max-time 2
curl localhost:8080/slow --max-time 2
curl localhost:8080/slow --max-time 2
192.168.65.1 - GET /slow HTTP/1.1 status=499 backend=- upstream=172.18.0.2:8080 request_time=2.005 upstream_time=2.005
192.168.65.1 - GET /slow HTTP/1.1 status=499 backend=- upstream=172.18.0.4:8080 request_time=2.005 upstream_time=2.00
```

### Mimic server crashes
```
curl localhost:8080/crash

This should result in the below
2025/12/26 17:32:30 [error] 29#29: *25 upstream prematurely closed connection while reading response header from upstream, client: 192.168.65.1, server: , request: "GET /crash HTTP/1.1", upstream: "http://172.18.0.3:8080/crash", host: "localhost:8080"
2025/12/26 17:32:30 [error] 29#29: *25 upstream prematurely closed connection while reading response header from upstream, client: 192.168.65.1, server: , request: "GET /crash HTTP/1.1", upstream: "http://172.18.0.2:8080/crash", host: "localhost:8080"
2025/12/26 17:32:30 [error] 29#29: *25 upstream prematurely closed connection while reading response header from upstream, client: 192.168.65.1, server: , request: "GET /crash HTTP/1.1", upstream: "http://172.18.0.4:8080/crash", host: "localhost:8080"
192.168.65.1 - GET /crash HTTP/1.1 status=502 backend=- upstream=172.18.0.3:8080, 172.18.0.2:8080, 172.18.0.4:8080 request_time=0.004 upstream_time=0.001, 0.001, 0.001
```
- This is normal - when one server fails, nginx will retry the next server
- This is why ideompotency is important
- To control retry behavior
    - `proxy_next_upstream off;` -> disables retries altogether
    - `proxy_next_upstream error timeout; proxy_next_upstream_tries 1;` -> safer, retries once

