package util

import (
	"fmt"
	"net"
)

var localIP string

// GetLocalIP 获取本地iP
func GetLocalIP() (ip string, err error) {
	if len(localIP) != 0 {
		ip = localIP
		return
	}

	netints, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, netint := range netints {
		if (netint.Flags & net.FlagUp) != 0 {
			addrs, addrErr := netint.Addrs()
			if addrErr != nil {
				err = addrErr
				return
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						localIP = ipnet.IP.String()
						ip = localIP
						return
					}
				}
			}
		}
	}
	err = fmt.Errorf("获取本地ip地址失败")
	return
}
