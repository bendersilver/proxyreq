package proxyreq

import (
	"net/http"
	"net/url"
	"time"

	"github.com/imroc/req"
	"golang.org/x/net/proxy"
)

// New -
func New(proxyHostPort, proxyType string) (*req.Req, error) {
	var dialer proxy.Dialer
	var err error
	var r *req.Req

	r = req.New()
	r.Client().Timeout = time.Duration(ClientTimeout) * time.Second

	ur := &url.URL{
		Host:   proxyHostPort,
		Scheme: proxyType,
	}
	if proxyType == "https" {
		dialer, err = proxy.FromURL(ur, directHTTPS{})
	} else {
		dialer, err = proxy.FromURL(ur, direct{})
	}
	if err == nil {
		tr := &http.Transport{
			Dial: dialer.Dial,
		}
		r.Client().Transport = tr
	}
	return r, err
}
