package proxyreq

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/textproto"
	"net/url"
	"strings"

	"golang.org/x/net/proxy"
)

type direct struct{}

func (direct) Dial(network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout: DialTimeout,
	}
	conn, err := dialer.Dial(network, addr)
	if err != nil {
		return conn, err
	}
	return conn, err
}

type directHTTPS struct{}

func (directHTTPS) Dial(network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout: DialTimeout,
	}
	return tls.DialWithDialer(dialer, network, addr, &tls.Config{})
}

// HTTP -
type HTTP struct {
	dialer proxy.Dialer
	addr   string
}

// newHTTP -
func newHTTP(host string, d proxy.Dialer) (*HTTP, error) {
	h := &HTTP{
		dialer: d,
		addr:   host,
	}
	return h, nil
}

// newHTTPDialer returns a http proxy dialer.
func newHTTPDialer(uri *url.URL, d proxy.Dialer) (proxy.Dialer, error) {
	return newHTTP(uri.Host, d)
}

func parseStartLine(line string) (r1, r2, r3 string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

// conn is a base conn struct.
type conn struct {
	r *bufio.Reader
	net.Conn
}

// Reader returns the internal bufio.Reader.
func (c *conn) Reader() *bufio.Reader {
	return c.r
}

// newConn returns a new conn.
func newConn(c net.Conn) *conn {
	return &conn{bufio.NewReader(c), c}
}

// Dial -
func (h *HTTP) Dial(network, addr string) (net.Conn, error) {
	c, err := h.dialer.Dial(network, h.addr)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.WriteString("CONNECT " + addr + " HTTP/1.1\r\n")
	buf.WriteString("Host: " + addr + "\r\n")
	buf.WriteString("Proxy-Connection: Keep-Alive\r\n")
	buf.WriteString("\r\n")

	_, err = c.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}

	con := newConn(c)
	tpr := textproto.NewReader(con.Reader())
	line, err := tpr.ReadLine()
	if err != nil {
		return con, err
	}
	_, code, _, ok := parseStartLine(line)
	if ok && code == "200" {
		tpr.ReadMIMEHeader()
		return con, err
	}
	switch code {
	case "403":
		return nil, fmt.Errorf("[http] 'CONNECT' to ports other than 443 are not allowed by proxy %s", h.addr)
	case "405":
		return nil, fmt.Errorf("[http] 'CONNECT' method not allowed by proxy %s", h.addr)
	case "407":
		return nil, fmt.Errorf("[http] authencation needed by proxy %s", h.addr)
	}

	return nil, fmt.Errorf("[http] can not connect remote address: %s. error code: %s", addr, code)

}

func init() {
	proxy.RegisterDialerType("http", newHTTPDialer)
	proxy.RegisterDialerType("https", newHTTPDialer)
}
