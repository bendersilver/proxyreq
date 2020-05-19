package proxyreq

import (
	"fmt"
	"testing"

	"github.com/imroc/req"
)

func reqCheck(r *req.Req) error {
	var ipify struct {
		IP *string `json:"ip"`
	}
	rsp, err := r.Get("https://api.ipify.org?format=json")
	if err != nil {
		return err
	}
	rsp.ToJSON(&ipify)
	// panic(ipify)
	if ipify.IP == nil {
		return fmt.Errorf("no parse ip")
	}
	return nil
}

func TestReqHTTPS(t *testing.T) {
	r, err := NewReq("us34.tcdn.me:443", "https")
	if err != nil {
		t.Fatal(err)
	}
	err = reqCheck(r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReqSocks5(t *testing.T) {
	r, err := NewReq("209.216.137.197:3129", "socks5")
	if err != nil {
		t.Fatal(err)
	}
	err = reqCheck(r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReqHttp(t *testing.T) {
	r, err := NewReq("190.144.34.146:3128", "http")
	if err != nil {
		t.Fatal(err)
	}
	err = reqCheck(r)
	if err != nil {
		t.Fatal(err)
	}
}
