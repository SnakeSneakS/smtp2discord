package main

import "net"

func IsClosedNetworkError(err error) bool {
	netErr, ok := err.(*net.OpError)
	return ok && netErr.Op == "accept"
}
