package util

import (
	"fmt"
	"net"
)

// GetIPFromCIDR Naive implementation, but works
func GetIPFromCIDR(cidr string, n int) (net.IP, error) {
	if n < 0 {
		return nil, fmt.Errorf("cannot compute n=%d IP address", n)
	}
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("error parsing CIDR: %v", err)
	}

	// Need a copy of the IP byte slice to reuse ipnet
	ip := make([]byte, len(ipnet.IP), cap(ipnet.IP))
	copy(ip, ipnet.IP)

	// Increment the IP address n times
	for j := 0; j < n; j++ {
		for i := len(ip) - 1; i >= 0; i-- {
			ip[i]++
			if ip[i] != 0 {
				break
			}
		}
	}

	// Verify the resulting IP is contained in the CIDR
	if !ipnet.Contains(ip) {
		return nil, fmt.Errorf("Could not compute the n=%d IP address of CIDR %q (resulting IP %q is not in CIDR)", n, cidr, net.IP(ip))
	}

	return ip, nil
}
