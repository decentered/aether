package api

import (
	"aether-core/services/globals"
	"aether-core/services/logging"
	"aether-core/services/tcpmim"
	"aether-core/services/toolbox"
	"fmt"
	"net"
	"time"
)

/*
The way we open a reverse connection is that we open a conn to the local server, and we open a conn to the remote server, and after sending the TCPMim message to request reverse open, we pipe one conn to another.

							C1	(connToLocal)								C2 (connToRemote)
LOCAL SERVER <--> LOCAL END <PIPE> LOCAL END <--> REMOTE SERVER
^ Local Remote    ^ Local Local    ^ Local Local  ^ Remote Remote

Why? Because it allows us to inject data to be sent to the remote server (the reverse open request) and *then* set up the pipe . It also saves us from having to deal with SO_REUSEADDR and SO_REUSEPORT and their subtly different platform-specific implementations.
*/

func RequestInboundSync(host string, subhost string, port uint16) {
	logging.Logf(1, "Attempting to request inbound sync from remote: %s/%s:%v", host, subhost, port)
	to := fmt.Sprint(host, ":", port)
	connToRemote, err := net.Dial("tcp4", to)
	if err != nil {
		logging.Logf(1, "Request inbound sync failed while attempting to establish a connection to the remote. Error: %v", err)
		return
	}
	localSrvAddr := fmt.Sprint(":", globals.BackendConfig.GetExternalPort())
	connToLocal, err := net.Dial("tcp4", localSrvAddr)
	if err != nil {
		logging.Logf(1, "Request inbound sync failed while attempting to establish a connection to the local server. Error: %v", err)
		return
	}
	// Set the values to transient config so that the server will be able to check if an incoming conn is a reverse conn.
	c1LocalLocalAddr, c1LocalLocalPort := toolbox.SplitHostPort(connToLocal.LocalAddr().String())
	globals.BackendTransientConfig.ReverseConnData.C1LocalLocalAddr = c1LocalLocalAddr
	globals.BackendTransientConfig.ReverseConnData.C1LocalLocalPort = c1LocalLocalPort
	mimMsg := tcpmim.MakeMimMessage(tcpmim.ReverseOpenRequest)
	fmt.Fprintf(connToRemote, string(mimMsg))
	// fmt.Fprintf(connToRemote, "YO\n")
	logging.Logf(1, "Established pipe: (Local End) R: %v -> L: %v >[Pipe]> R: %v > L: %v (Remote End)",
		connToLocal.RemoteAddr().String(),
		connToLocal.LocalAddr().String(),
		connToRemote.LocalAddr().String(),
		connToRemote.RemoteAddr().String(),
	)
	start := time.Now()
	// Set timeouts to infinite - both are successful.
	pipe(connToRemote, connToLocal)
	// The remote will auto-close the connection, or the local server will, or it will just timeout on its own based on inactivity.
	elapsed := time.Since(start)
	fmt.Printf("reverse conn took %v\n", elapsed)
}

func connToChan(conn net.Conn) chan []byte {
	c := make(chan []byte)
	go func() {
		b := make([]byte, 1024)
		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				// Copy just so so it doesn't change while read by the recipient
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()
	return c
}

func pipe(c1, c2 net.Conn) {
	chan1 := connToChan(c1)
	chan2 := connToChan(c2)
	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				return
			} else {
				// In the case of a r/w, update expiries.
				updateDeadlines(&c1, &c2)
				// c1.SetDeadline(time.Now().Add(1 * time.Minute))
				// c2.SetDeadline(time.Now().Add(1 * time.Minute))
				c2.Write(b1)
			}
		case b2 := <-chan2:
			if b2 == nil {
				return
			} else {
				// In the case of a r/w, update expiries.
				updateDeadlines(&c1, &c2)
				// c1.SetDeadline(time.Now().Add(1 * time.Minute))
				// c2.SetDeadline(time.Now().Add(1 * time.Minute))
				c1.Write(b2)
			}
		}
	}
}

var lastDeadlineUpdate int64

func updateDeadlines(c1, c2 *net.Conn) {
	if lastDeadlineUpdate < time.Now().Add(-30*time.Second).Unix() {
		(*c1).SetDeadline(time.Now().Add(1 * time.Minute))
		(*c2).SetDeadline(time.Now().Add(1 * time.Minute))
	}
}
