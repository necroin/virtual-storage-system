package lan

import (
	"net"
	"strings"
	"vss/src/settings"
)

func GetMyLanAddr() string {
	result := "localhost" + settings.DefaultPort
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return result
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if strings.Contains(ip, "192.168.0.") {
					return ip + settings.DefaultPort
				}
			}
		}

	}

	return result
}
