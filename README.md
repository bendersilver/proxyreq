# proxyreq

A golang http proxy request library based from https://github.com/imroc/req

# settings

DialTimeout (default 3s)

ClientTimeout (default 3s)

### Example

```go
package main

import (
	"fmt"

	"github.com/bendersilver/proxyreq"
)

func main() {
	var proxy = "209.216.137.197:3129" // host:port
	var proxyType = "socks5"           // http | https | socks5
	r, err := proxyreq.New(proxy, proxyType)
	if err != nil {
		panic(err)
	}
	// add head
	// "User-Agent":      browser.Computer()
	// "Accept-Encoding": "gzip"
	rsp, err := r.Get("https://api.ipify.org")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", rsp.String())
	// out: proxy ip
}
```