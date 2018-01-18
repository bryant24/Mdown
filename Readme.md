## Introduce
Golang multi routine downloader
1、Support multi-routine download  <br />
2、Break point and continuous download  <br />
3、Download error retry(unexpected EOF,etc)  <br />


## Install
```
go get -u https://github.com/bryant24/Mdown
```

## How to use
```
package main

import (
	"strings"
	"github.com/bryant24/Mdown"
)

func main() {
	src := "https://abc.com/abc.zip"
	wg:=&sync.WaitGroup{}  // lock
    timeout:=30*time.Second // timeout for file chunk

	Mdown.MultiDownload(url, "a.zip", timeout,wg)
}

```
