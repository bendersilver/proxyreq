package proxyreq

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	"github.com/bendersilver/proxyreq/dhttp"
	"github.com/bendersilver/proxyreq/dhttps"
	"golang.org/x/net/proxy"
)

// host -
type host struct {
	forward proxy.Dialer
	addr    string
}

// newHost -
func newHost(hst string, d proxy.Dialer) (*host, error) {
	h := &host{
		forward: d,
		addr:    hst,
	}
	return h, nil
}

// newHostDialer returns a http proxy dialer.
func newHostDialer(uri *url.URL, d proxy.Dialer) (proxy.Dialer, error) {
	return newHost(uri.Host, d)
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
func (h *host) Dial(network, addr string) (net.Conn, error) {
	c, err := h.forward.Dial(network, h.addr)
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

// Dialer -
func Dialer(h, t string) (proxy.Dialer, error) {
	ur := &url.URL{
		Host:   h,
		Scheme: t,
	}
	switch t {
	case "https":
		return proxy.FromURL(ur, dhttps.Direct)
	case "http":
		return proxy.FromURL(ur, dhttp.Direct)
	case "socks5":
		return proxy.SOCKS5("tcp", h, nil, dhttp.Direct)
	}
	return nil, fmt.Errorf("err")
}

// SetDialTimeout -
func SetDialTimeout(t time.Duration) {
	dhttp.DialTimeout = t
	dhttps.DialTimeout = t
}

// SetConnTimeout -
func SetConnTimeout(t time.Duration) {
	dhttp.ConnTimeout = t
	dhttps.ConnTimeout = t
}

func init() {
	proxy.RegisterDialerType("http", newHostDialer)
	proxy.RegisterDialerType("https", newHostDialer)
}
