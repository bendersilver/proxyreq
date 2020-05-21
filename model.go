package proxyreq

import "github.com/imroc/req"

// Rq -
type Rq struct {
	rq *req.Req
}

// Post -
func (r *Rq) Post(url string, args ...interface{}) (*req.Resp, error) {
	return Post(r.rq, url, args...)
}

// Get -
func (r *Rq) Get(url string, args ...interface{}) (*req.Resp, error) {
	return Get(r.rq, url, args...)
}

// NewModel -
func NewModel(proxyHostPort, proxyType string) (*Rq, error) {
	r, err := New(proxyHostPort, proxyType)
	if err != nil {
		return nil, err
	}
	return &Rq{
		rq: r,
	}, nil
}
