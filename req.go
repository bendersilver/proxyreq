package proxyreq

import (
	"crypto/tls"
	"net/http"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/imroc/req"
	"golang.org/x/net/proxy"
)

func getArgs(args ...interface{}) []interface{} {
	args = append(args, req.Header{
		"User-Agent":      browser.Computer(),
		"Accept-Encoding": "gzip",
	})
	args = append(args)
	return args
}

// post -
func post(rq *req.Req, url string, args ...interface{}) (*req.Resp, error) {
	return rq.Post(url, getArgs(args)...)
}

// get -
func get(rq *req.Req, url string, args ...interface{}) (*req.Resp, error) {
	return rq.Get(url, getArgs(args)...)
}

// new -
func new(proxyHostPort, proxyType string) (*req.Req, error) {
	var dialer proxy.Dialer
	var err error
	var r *req.Req

	r = req.New()
	dialer, err = Dialer(proxyHostPort, proxyType)
	if err == nil {
		r.Client().Transport = &http.Transport{
			Dial: dialer.Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return r, err
}

// NewEmpty -
func NewEmpty() *Rq {
	return &Rq{
		rq: req.New(),
	}
}

// SetTransport -
func (r *Rq) SetTransport(proxyHostPort, proxyType string) error {
	dialer, err := Dialer(proxyHostPort, proxyType)
	if err == nil {
		r.rq.Client().Transport = &http.Transport{
			Dial: dialer.Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return err
}
