package dhttps

import (
	"crypto/tls"
	"net"
	"time"
)

// DialTimeout -
var DialTimeout = time.Second * 3

// ConnTimeout -
var ConnTimeout = time.Second * 3

type direct struct{}

// Direct implements Dialer by making network connections directly using net.Dial or net.DialContext.
var Direct = direct{}

// Dial directly invokes net.Dial with the supplied parameters.
func (direct) Dial(network, addr string) (conn net.Conn, err error) {
	var d net.Dialer
	d.Timeout = DialTimeout
	conn, err = tls.DialWithDialer(&d, network, addr, &tls.Config{})
	if err != nil {
		return
	}
	conn.SetDeadline(time.Now().Add(ConnTimeout))
	return
}
