package addr

import (
	"net"
	"net/http"
	"strings"
)

type Option func(*Options)

type Options struct {
	FlagUp       bool
	FlagLoopback bool

	OnlyV4 bool
}

func NetFlagUP(flag bool) Option {
	return func(o *Options) {
		o.FlagUp = flag
	}
}

func NetFlagLoopback(flag bool) Option {
	return func(o *Options) {
		o.FlagLoopback = flag
	}
}

func NetOnlyV4(flag bool) Option {
	return func(o *Options) {
		o.OnlyV4 = flag
	}
}

func ExternalIPs() []string {
	return IPs(NetOnlyV4(true), NetFlagUP(true), NetFlagLoopback(true))
}

func IPs(ofs ...Option) []string {
	opts := &Options{}

	for _, o := range ofs {
		o(opts)
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ips []string
	for _, iface := range ifaces {
		if opts.FlagUp && iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if opts.FlagLoopback && iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := GetIPFromAddr(addr, opts.OnlyV4)
			if ip == nil {
				continue
			}

			ips = append(ips, ip.String())
		}
	}

	return ips
}

// GetIPFromAddr get ip from addr
func GetIPFromAddr(addr net.Addr, onlyV4 bool) net.IP {
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
	if onlyV4 {
		ip = ip.To4()
		if ip == nil {
			return nil // not an ipv4 address
		}
	}

	return ip
}

// GetClientIP get client ip from http request
func GetClientIP(ctx *http.Request) string {
	if ctx == nil {
		return ""
	}

	// Cdn-Src-Ip
	if ip := ctx.Header.Get("Cdn-Src-Ip"); ip != "" {
		return ip
	}

	// X-Forwarded-For
	if ips := ctx.Header.Get("X-Forwarded-For"); ips != "" {
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
	if ip := ctx.Header.Get("Client-Ip"); ip != "" {
		return ip
	}

	// RemoteAddr
	if ip, _, err := net.SplitHostPort(ctx.RemoteAddr); err == nil {
		return ip
	}

	return ""
}
