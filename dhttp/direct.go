package dhttp

import (
	"context"
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
	conn, err = d.Dial(network, addr)
	if err != nil {
		return
	}
	conn.SetDeadline(time.Now().Add(ConnTimeout))
	return
}

// DialContext instantiates a net.Dialer and invokes its DialContext receiver with the supplied parameters.
func (direct) DialContext(ctx context.Context, network, addr string) (conn net.Conn, err error) {
	var d net.Dialer
	d.Timeout = DialTimeout
	conn, err = d.DialContext(ctx, network, addr)
	if err != nil {
		return
	}
	conn.SetDeadline(time.Now().Add(ConnTimeout))
	return conn, err
}
