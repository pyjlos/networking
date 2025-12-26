## Simple http server

## Endpoints
```
/fast
/slow
/crash
```

## Start server
```
go run main.go
```

## Client calls (using curl)
1. Basic http call to each
```
curl localhost:8080/fast
curl localhost:8080/slow -> after 5 seconds should response
curl localhost:8080/crash -> should crash
```

2. Client timeouts
```
curl localhost:8080/fast --max-time 2
curl localhost:8080/slow --max-time 2 -> should return curl: (28) Operation timed out after 2002 milliseconds with 0 bytes received
curl localhost:8080/slow -> shut down server -> should return curl: (52) Empty reply from server
```
- The empty reply from the server means there was no response
- The operation timeout -> means it was able to connect to the server, but it got no response back

3. Client Connection Keep-alive vs. close
```
curl -H "Connection: keep-alive" localhost:8080/fast
    - keep-alive is default
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /fast HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> Connection: keep-alive
>
* Request completely sent off
< HTTP/1.1 200 OK
< Date: Fri, 26 Dec 2025 16:09:22 GMT
< Content-Length: 26
< Content-Type: text/plain; charset=utf-8
<
Hello from fast endpoint!
* Connection #0 to host localhost left intact
```

```
curl -H "Connection: close" localhost:8080/fast
❯ curl -H "Connection: close" localhost:8080/fast -v
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /fast HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> Connection: close
>
* Request completely sent off
< HTTP/1.1 200 OK
< Date: Fri, 26 Dec 2025 16:10:32 GMT
< Content-Length: 26
< Content-Type: text/plain; charset=utf-8
< Connection: close
<
Hello from fast endpoint!
* Closing connection
```
- generally speaking -> keep-alive is probably what we want
- only use close if you know you only need to make one call