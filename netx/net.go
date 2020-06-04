package helper

import (
	"net"
)

func IpAddrNoPort(s string) string {
	ip, _, err := net.SplitHostPort(s)
	if err == nil {
		return ip
	}

	ip2 := net.ParseIP(s)
	if ip2 == nil {
		return ""
	}

	return ip2.String()
}
