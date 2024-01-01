package lan

import (
	"net"
	"strings"
	"vss/src/settings"
)

func GetMyLanAddr() string {
	result := "localhost" + settings.DefaultPort
	ifaces, err := net.Interfaces()
	if err != nil {
		return result
	}
	for _, i := range ifaces {
		if strings.Contains(i.Flags.String(), "up") {
			addrs, err := i.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				ip := ""
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP.String()
				case *net.IPAddr:
					ip = v.IP.String()
				}
				if strings.Contains(ip, "192.168.") {
					return ip + settings.DefaultPort
				}
			}
		}
	}

	return result
}
