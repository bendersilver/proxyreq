# proxyreq

A golang http proxy request library based from https://github.com/imroc/req

# settings

DialTimeout (default 3s)

ClientTimeout (default 3s)

### Example

```go
package main

import "github.com/bendersilver/proxyreq"


func main() {
    var proxy = "209.216.137.197:3129" // host:port
    var proxyType = "socks5" // http | https | socks5
    r, err := proxyreq.New(proxy, proxyType)
    if err != nil {
        panic(err)
    }
    rsp, err := r.Get("https://api.ipify.org?format=json")
    if err != nil {
        panic(err)
    }
	print(rsp.String())
}


```