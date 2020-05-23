package proxyreq

import (
	"net/http"
	"net/url"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/bendersilver/proxyreq/dhttp"
	"github.com/bendersilver/proxyreq/dhttps"
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
func (r *Rq) SetTransport(u *url.URL) error {
	var dialer proxy.Dialer
	var err error
	if u.Scheme == "https" {
		dialer, err = proxy.FromURL(u, dhttps.Direct)
	} else {
		dialer, err = proxy.FromURL(u, dhttp.Direct)
	}
	r.rq.Client().Transport = &http.Transport{
		Dial: dialer.Dial,
	}
	return err
}
