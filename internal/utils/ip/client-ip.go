package ip

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

type IPOutput struct {
	RealIP     string
	Port       string
	ProxyChain []string
}

func ReadUserIP(r *http.Request) (*IPOutput, error) {
	output := &IPOutput{}

	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		proxyChain, err := parseAndValidateIPs(xForwardedFor)
		if err != nil {
			return nil, err
		}
		output.ProxyChain = proxyChain
		output.RealIP = proxyChain[0]

		_, port, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			output.Port = port
		}
		return output, nil
	}

	ip := r.Header.Get("X-Real-Ip")
	if ip != "" {
		if !isValidIP(ip) {
			return nil, fmt.Errorf("invalid IP address: %s", ip)
		}
		output.RealIP = ip

		_, port, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			output.Port = port
		}
		return output, nil
	}

	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		output.RealIP = r.RemoteAddr
		return output, nil
	}
	output.RealIP = host
	output.Port = port

	return output, nil
}

func isValidIP(ip string) bool {
	if strings.Contains(ip, "%") {
		ip = strings.Split(ip, "%")[0]
	}
	return net.ParseIP(ip) != nil
}

func parseAndValidateIPs(ips string) ([]string, error) {
	ipList := strings.Split(ips, ",")
	var validIPs []string

	for _, ip := range ipList {
		ip = strings.TrimSpace(ip)
		if strings.Contains(ip, "%") {
			ip = strings.Split(ip, "%")[0]
		}
		if !isValidIP(ip) {
			return nil, fmt.Errorf("invalid IP address in proxy chain: %s", ip)
		}
		validIPs = append(validIPs, ip)
	}

	return validIPs, nil
}
