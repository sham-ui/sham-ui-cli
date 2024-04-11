package net_addr

import (
	"cms/pkg/asserts"
	"testing"
)

func TestResolve(t *testing.T) {
	testCases := []struct {
		addr                   string
		expectedNetwork        string
		expectedClearedAddress string
	}{
		{
			addr:                   "unix:./foo.sock",
			expectedNetwork:        UnixNetwork,
			expectedClearedAddress: "./foo.sock",
		},
		{
			addr:                   "/tmp/foo.sock",
			expectedNetwork:        UnixNetwork,
			expectedClearedAddress: "/tmp/foo.sock",
		},
		{
			addr:                   "localhost:80",
			expectedNetwork:        TcpNetwork,
			expectedClearedAddress: "localhost:80",
		},
		{
			addr:                   ":80",
			expectedNetwork:        TcpNetwork,
			expectedClearedAddress: ":80",
		},
	}
	for _, test := range testCases {
		t.Run(test.addr, func(t *testing.T) {
			// Action
			network, addr := Resolve(test.addr)

			// Assert
			asserts.Equals(t, test.expectedNetwork, network)
			asserts.Equals(t, test.expectedClearedAddress, addr)
		})
	}
}
