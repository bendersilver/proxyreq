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
	if ipify.IP == nil {
		return fmt.Errorf("no parse ip")
	}
	return nil
}

func TestReqHTTPS(t *testing.T) {
	r, err := New("us34.tcdn.me:443", "https")
	if err != nil {
		t.Fatal(err)
	}
	err = reqCheck(r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReqSocks5(t *testing.T) {
	r, err := New("209.216.137.197:3129", "socks5")
	if err != nil {
		t.Fatal(err)
	}
	err = reqCheck(r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReqHttp(t *testing.T) {
	r, err := New("190.144.34.146:3128", "http")
	if err != nil {
		t.Fatal(err)
	}
	err = reqCheck(r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReqGet(t *testing.T) {
	r, err := New("190.144.34.146:3128", "http")
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := Get(r, "https://www.avito.ru/sverdlovskaya_oblast/oborudovanie_dlya_biznesa?s=104&user=1&q=%D1%81%D1%82%D1%83%D0%BB+%D0%B4%D0%BB%D1%8F+%D0%BE%D1%84%D0%B8%D1%81%D0%B0")
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(rsp.String())
	var ipify struct {
		IP *string `json:"ip"`
	}
	rsp.ToJSON(&ipify)
	if ipify.IP == nil {
		t.Fatal("no parse ip")
	}
	t.Fatal(rsp.String())
}
