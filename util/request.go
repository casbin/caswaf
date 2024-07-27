package util

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIp(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		clientIP := strings.Split(forwarded, ",")[0]
		return strings.TrimSpace(clientIP)
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
