// Services > Ports
// This package provides function related to ports in the local machine, such as finding a free port, or checking whether a port is open for use.

package ports

import (
	"aether-core/services/globals"
	"aether-core/services/logging"
	"fmt"
	"net"
	"strings"
)

// GetFreePort returns a free port that is currently unused in the local system.
func GetFreePort() int {
	a, err := net.ResolveTCPAddr("tcp4", ":0")
	if err != nil {
		logging.LogCrash(fmt.Sprintf("We could not parse the TCP address in an attempt to get a free port. The error raised was: %s", err))
	}
	l, err := net.ListenTCP("tcp4", a)
	defer l.Close()
	if err != nil {
		logging.LogCrash(fmt.Sprintf("We could not listen to TCP in an attempt to get a free port. The error raised was: %s", err))
	}
	return l.Addr().(*net.TCPAddr).Port
}

// GetFreePorts returns a number of free ports that are currently unused in the local system.
func GetFreePorts(number int) []int {
	ports := []int{}
	clashcount := 0
	checkport := func(ports []int, port int) bool {
		for _, val := range ports {
			if val == port {
				clashcount++
				if clashcount > 65535 {
					logging.LogCrash(fmt.Sprintf("This computer does not have enough ports that are free. You've requested %d free ports. ", number))
				}
				return true
			}
		}
		return false
	}
	for i := 0; i < number; i++ {
		port := GetFreePort()
		for checkport(ports, port) {
			port = GetFreePort()
		}
		ports = append(ports, port)
	}
	return ports
}

// CheckPortAvailability checks for whether a port that it is given is currently free to use.
func CheckPortAvailability(port uint16) bool {
	a, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		logging.LogCrash(fmt.Sprintf("We could not parse the TCP address in an attempt to check the availability of the given port. The error raised was: %s, The port attempted to be checked was: %d", err, port))
	}
	l, err := net.ListenTCP("tcp4", a)
	defer l.Close()
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") || strings.Contains(err.Error(), "permission denied") {
			logging.Log(1, fmt.Sprintf("Port number %d is already in use. Error: %s", port, err))
			return false
		} else {
			logging.LogCrash(fmt.Sprintf("We attempted to check the availability of the port %d on the current computer and it failed with this error:", err, port))
		}
	}
	return true
}

// VerifyLocalPort verifies the local port available in the config, and if it is not available, replaces it with one that is. Then it flips the bit to mark the local port as verified.
func VerifyExternalPort() {
  logging.Log(2,"VerifyExternalPort check is running.")
	// This check only runs once per start.
	if !globals.BackendTransientConfig.ExternalPortVerified {
		// Prevent race condition in which any number of calls can enter this before CheckPortAvailability returns.
		globals.BackendTransientConfig.ExternalPortVerified = true
		if CheckPortAvailability(globals.BackendConfig.GetExternalPort()) {
			logging.Log(1, fmt.Sprintf("The external port %d is verified to be open and available for use.", globals.BackendConfig.GetExternalPort()))
		} else {
			freeport := GetFreePort()
			logging.Log(1, fmt.Sprintf("The port number for this node has changed. New external port for this node is: %d", freeport))
			globals.BackendConfig.SetExternalPort(freeport)
		}
	}
  logging.Log(2,"VerifyExternalPort check is done.")
}
