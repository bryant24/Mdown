## Introduce

Golang multi routine downloader   <br/>
1、Support multi-goroutine download  <br />
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
	origin := "https://abc.com/abc.zip"
	target := "a.zip"
        timeout:= 30
	downloader := NewMDownloader(origin, target, timeout)
        downloader.Start()
}

```
