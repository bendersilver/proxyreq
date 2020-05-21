package proxyreq

import (
	"context"
	"net/http"
	"net/url"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/imroc/req"
	"golang.org/x/net/proxy"
)

func getArgs(args ...interface{}) (*time.Timer, []interface{}) {
	args = append(args, req.Header{
		"User-Agent":      browser.Computer(),
		"Accept-Encoding": "gzip",
	})
	ctx, cancel := context.WithCancel(context.TODO())
	timer := time.AfterFunc(ClientTimeout, func() {
		cancel()
	})
	args = append(args, ctx)
	return timer, args
}

// post -
func post(rq *req.Req, url string, args ...interface{}) (*req.Resp, error) {
	timer, v := getArgs(args)
	defer timer.Stop()
	return rq.Post(url, v...)
}

// get -
func get(rq *req.Req, url string, args ...interface{}) (*req.Resp, error) {
	timer, v := getArgs(args)
	defer timer.Stop()
	return rq.Get(url, v...)
}

// new -
func new(proxyHostPort, proxyType string) (*req.Req, error) {
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
