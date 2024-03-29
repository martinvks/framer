package utils

import (
	"fmt"
	"net"
)

func LookUp(domain string) (net.IP, error) {
	ips, err := net.LookupIP(domain)

	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4, nil
		}
	}

	return nil, fmt.Errorf("DNS lookup failed for %v", domain)
}
