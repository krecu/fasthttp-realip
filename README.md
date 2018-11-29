# RealIP [fasthttp]

[![GoDoc](https://godoc.org/github.com/krecu/fasthttp-realip?status.svg)](http://godoc.org/github.com/krecu/fasthttp-realip)

Go package that can be used to get client's real public IP, which usually useful for logging HTTP server.

### Feature

* Follows the rule of X-Real-IP
* Follows the rule of X-Forwarded-For
* Exclude local or private address

## Example

```go
package main

import "github.com/krecu/fasthttp-realip"

func RealIpHandler(ctx *fasthttp.RequestCtx) {
	clientIP := realip.FromRequest(ctx)
	log.Println("GET / from", clientIP)
}
```

## Developing

Commited code must pass:

* [golint](https://github.com/golang/lint)
* [go vet](https://godoc.org/golang.org/x/tools/cmd/vet)
* [gofmt](https://golang.org/cmd/gofmt)
* [go test](https://golang.org/cmd/go/#hdr-Test_packages):
