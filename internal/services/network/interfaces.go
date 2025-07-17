package network

import (
	"net"
)

func GetInterfaces() []string {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var result []string
	for _, iface := range ifs {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 && len(iface.HardwareAddr) > 0 {
			result = append(result, iface.Name+": "+iface.HardwareAddr.String())
		}
	}
	return result
}
