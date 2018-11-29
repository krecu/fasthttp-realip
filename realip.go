package realip

import (
	"errors"
	"net"
	"strings"

	"github.com/valyala/fasthttp"
)

var cidrs []*net.IPNet

func init() {
	maxCidrBlocks := []string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}

	cidrs = make([]*net.IPNet, len(maxCidrBlocks))
	for i, maxCidrBlock := range maxCidrBlocks {
		_, cidr, _ := net.ParseCIDR(maxCidrBlock)
		cidrs[i] = cidr
	}
}

// isLocalAddress works by checking if the address is under private CIDR blocks.
// List of private CIDR blocks can be seen on :
//
// https://en.wikipedia.org/wiki/Private_network
//
// https://en.wikipedia.org/wiki/Link-local_address
func isPrivateAddress(address string) (bool, error) {
	ipAddress := net.ParseIP(address)
	if ipAddress == nil {
		return false, errors.New("address is not valid")
	}

	for i := range cidrs {
		if cidrs[i].Contains(ipAddress) {
			return true, nil
		}
	}

	return false, nil
}

// FromRequest return client's real public IP address from http request headers.
func FromRequest(ctx *fasthttp.RequestCtx) string {
	// Fetch header value
	xRealIP := ctx.Request.Header.Peek("X-Real-Ip")
	xForwardedFor := ctx.Request.Header.Peek("X-Forwarded-For")

	// If both empty, return IP from remote address
	if xRealIP == nil && xForwardedFor == nil {
		var remoteIP string

		// If there are colon in remote address, remove the port number
		// otherwise, return remote address as is
		RemoteAddr := ctx.RemoteAddr().String()
		if strings.ContainsRune(RemoteAddr, ':') {
			remoteIP, _, _ = net.SplitHostPort(RemoteAddr)
		} else {
			remoteIP = RemoteAddr
		}

		return remoteIP
	}

	// Check list of IP in X-Forwarded-For and return the first global address
	for _, address := range strings.Split(string(xForwardedFor), ",") {
		address = strings.TrimSpace(address)
		isPrivate, err := isPrivateAddress(address)
		if !isPrivate && err == nil {
			return address
		}
	}

	// If nothing succeed, return X-Real-IP
	return string(xRealIP)
}

// RealIP is depreciated, use FromRequest instead
func RealIP(ctx *fasthttp.RequestCtx) string {
	return FromRequest(ctx)
}
