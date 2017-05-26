package library

import (
	"net"
	"bytes"
	"../objects"
	log "../../core/logger"
)

var ip = net.ParseIP(objects.CoreIP)

// Parse the IP and see if they are in the range
func checkIP(found_ip string) bool {

	// Parse the IP
	obtain_ip := net.ParseIP(found_ip)

	// Get the local IP by comparing it.
	if bytes.Compare(obtain_ip, ip) >= 0 {
		return true
	}

	 return false
}

// Get the IP address
func GetLocalIP() (string, error){

	log.Println("Getting local IP address")

	// Get Interface address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", nil
	}

	// Extract the IP's and see if its in the range.
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if checkIP(ipnet.IP.String()) {
					log.Println("Found Local IP address: " + ipnet.IP.String())
					return ipnet.IP.String(), nil
				}
			}
		}
	}

	// If nothing found then produce a warning message
	log.Warn("Cannot find local IP address in the range \"" + objects.CoreIP + "\", ignoring the IP....")

	return "", nil
}