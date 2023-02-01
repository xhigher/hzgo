package env

import (
	"errors"
	"net"
	"os"
	"strings"
)

const (
	Local = "local"
	Dev   = "dev"
	Test  = "test"
	Prod  = "prod"
)

var env = Test

func Init(e string) {
	if strings.TrimSpace(e) == "" {
		env = Test
	}
	env = e
}

func IsLocal() bool { return env == Local }
func IsTest() bool  { return env == Test }
func IsDev() bool   { return env == Dev }
func IsProd() bool  { return env == Prod }

func externalIP() (net.IP, error) {
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
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("not connected to the network")
}

func getIpFromAddr(addr net.Addr) net.IP {
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

func GetLocalIP() string {
	ip, err := externalIP()
	if err != nil {
		return ""
	}
	return ip.String()
}

func GetHostName() string{
	name, err := os.Hostname()
	if err != nil {
		return ""
	}
	return name
}
