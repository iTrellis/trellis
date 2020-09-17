/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package internal

import (
	"errors"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := GetIPFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func GetIPFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func GetClientIP(ctx *gin.Context) string {

	// Cdn-Src-Ip
	if ip := ctx.GetHeader("Cdn-Src-Ip"); ip != "" {
		return ip
	}

	// X-Forwarded-For
	if ips := ctx.GetHeader("X-Forwarded-For"); ips != "" {
		addr := strings.Split(ips, ",")
		if len(addr) > 0 && addr[0] != "" {
			rip, _, err := net.SplitHostPort(addr[0])
			if err != nil {
				rip = addr[0]
			}
			return rip
		}
	}

	// Client_Ip
	if ip := ctx.GetHeader("Client-Ip"); ip != "" {
		return ip
	}

	// RemoteAddr
	if ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr); err == nil {
		return ip
	}

	return ""
}
