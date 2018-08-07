// Server > InspectTCP

// This is the Schrödinger's cat of libraries. It allows you to observe what is going on coming from the wire in a net.Conn without actually changing its state nor causing it to be affected in any way. Internally, what it does is that it actually caches the read and at the first read from the outside, the already-read part is provided out, so to the external world, the connection appears undisturbed and unobserved.

// The second part, the listener uses this ability to inspect the connection, and if it's something that we are interested in, stop propagation of that into the http Server, and assume direct control of that TCP connection.

package server

import (
	"aether-core/backend/dispatch"
	"aether-core/io/api"
	"aether-core/services/logging"
	"aether-core/services/tcpmim"
	// "aether-core/services/toolbox"
	"bufio"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	"net"
	// "strconv"
	// "io"
	// "io/ioutil"
	"time"
)

// max uint8: ROR MIM/255.65535

type InspectingListener struct {
	netListener net.Listener
}

func (l InspectingListener) Accept() (net.Conn, error) {
	// logging.Logf(1, "Accept inspector enters")
	nc, err := l.netListener.Accept()
	if err != nil {
		return nil, err
	}
	nic := NewInspectableConn(nc)
	data, _ := nic.Peek(5)
	dataStr := string(data)
	// logging.Logf(1, "This is what we got: %v", dataStr)
	if dataStr[0:3] == "MIM" { // This is a raw TCP Mim request.
		// We're setting up the deadlines, because if we read from malformed Mim response that lies in its length, the read will hang. This timeout means we'll eventually close the connection and move on in that case.
		deadline := time.Now().Add(60 * time.Second)
		nic.SetDeadline(deadline)
		nic.readDeadline = deadline
		len := int(uint8(dataStr[4:5][0]))
		msg, err := nic.Peek(len)
		if err != nil {
			panic(err)
		}
		mimMessage := tcpmim.ParseMimMessage(msg)
		if mimMessage == tcpmim.ReverseOpenRequest {
			logging.Logf(2, "This is a TCP reverse connect request.")
			if time.Now().Unix() < deadline.Unix() {
				// We've managed to come here without getting past the deadline. Let's reset the connection to have no deadline so that it won't cut off prematurely while we're doing what we want. 10 minutes is a last-resort guess in case something goes very wrong. Fetch() will set read deadlines as needed in the case it's dealing with a reverse conn.
				nc.SetDeadline(time.Now().Add(10 * time.Minute))
				dispatch.Sync(api.Address{}, []string{}, &nc)
			} else {
				logging.Logf(2, "Deadline exceeded while waiting for read. Passing this by to the server.")
			}
		}
		// This is where we determine if we want to reverse open into this. This is a high risk action.
	}
	// logging.Logf(1, "Accept inspector is done.")
	return nic, nil
}

func (l InspectingListener) Close() error {
	return l.netListener.Close()
}

func (l InspectingListener) Addr() net.Addr {
	return l.netListener.Addr()
}

type InspectableConn struct {
	r *bufio.Reader
	net.Conn
	readDeadline time.Time
}

func NewInspectableConn(c net.Conn) InspectableConn {
	return InspectableConn{bufio.NewReader(c), c, time.Time{}}
}

func NewInspectableConnSize(c net.Conn, n int) InspectableConn {
	return InspectableConn{bufio.NewReaderSize(c, n), c, time.Time{}}
}

func (ic InspectableConn) Peek(n int) ([]byte, error) {
	return ic.r.Peek(n)
}

func (ic InspectableConn) Read(p []byte) (int, error) {
	return ic.r.Read(p)
}

func (ic InspectableConn) Close() error {
	// logging.Logf(1, "Close was called")
	return ic.Conn.Close()
}

// func (ic InspectableConn) ReadBytes(p []byte) (int, error) {
// 	return ic.r.ReadBytes()
// }
