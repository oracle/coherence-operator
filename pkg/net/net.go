/*
 * Copyright (c) 2019, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package net

import "net"

// GetLocalAddress returns the local IP address.
func GetLocalAddress() (net.IP, error) {
	var flagP2P = net.FlagUp | net.FlagPointToPoint
	var mtuBest = 0
	var ipBest net.IP

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range ifaces {
		if (i.Flags&flagP2P) != 0 && i.MTU > mtuBest {
			addrs, err := i.Addrs()
			if err != nil {
				return nil, err
			}

			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP.To4()
				case *net.IPAddr:
					ip = v.IP.To4()
				}

				if ip != nil && !ip.IsLoopback() {
					mtuBest = i.MTU
					ipBest = ip
				}
			}
		}
	}

	return ipBest, nil
}
