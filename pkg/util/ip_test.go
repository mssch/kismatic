package util

import (
	"net"
	"testing"
)

func TestGetIPFromCIDR(t *testing.T) {
	tests := []struct {
		cidr       string
		n          int
		expectedIP net.IP
		expectErr  bool
	}{
		{
			cidr:       "10.5.6.217/32",
			n:          0,
			expectedIP: net.ParseIP("10.5.6.217"),
		},
		{
			cidr:      "10.5.6.217/32",
			n:         1,
			expectErr: true,
		},
		{
			// max addresses = 256
			cidr:       "10.20.0.0/24",
			n:          100,
			expectedIP: net.IP{byte(10), byte(20), byte(0), byte(100)},
		},
		{
			// max addresses = 256
			cidr:      "10.20.0.0/24",
			n:         300,
			expectErr: true,
		},
		{
			// max addresses = 16
			cidr:      "10.20.0.0/28",
			n:         16,
			expectErr: true,
		},
		{
			// max addresses = 16
			cidr:       "10.20.0.0/28",
			n:          15,
			expectedIP: net.ParseIP("10.20.0.15"),
		},
		{
			// max addresses = 16
			cidr:      "10.20.0.0/28",
			n:         -1,
			expectErr: true,
		},
		{
			cidr:       "172.16.0.0/16",
			n:          1,
			expectedIP: net.ParseIP("172.16.0.1"),
		},
		{
			cidr:       "172.16.0.0/16",
			n:          2,
			expectedIP: net.ParseIP("172.16.0.2"),
		},
	}

	for i, test := range tests {
		ip, err := GetIPFromCIDR(test.cidr, test.n)
		if err != nil {
			if !test.expectErr {
				t.Errorf("test %d - got an unexpected error: %v", i, err)
			}
			continue
		}

		if !ip.Equal(test.expectedIP) {
			t.Errorf("expected %q, but got %q", test.expectedIP, ip)
		}

		if test.expectErr {
			t.Errorf("expected an error, but didn't get one")
		}
	}
}
