package proxyreq

import "github.com/imroc/req"

// Rq -
type Rq struct {
	rq *req.Req
}

// Req -
func (r *Rq) Req() *req.Req {
	return r.rq
}

// Post -
func (r *Rq) Post(url string, args ...interface{}) (*req.Resp, error) {
	return post(r.rq, url, args...)
}

// Get -
func (r *Rq) Get(url string, args ...interface{}) (*req.Resp, error) {
	return get(r.rq, url, args...)
}

// New -
func New(proxyHostPort, proxyType string) (*Rq, error) {
	r, err := new(proxyHostPort, proxyType)
	if err != nil {
		return nil, err
	}
	return &Rq{
		rq: r,
	}, nil
}
