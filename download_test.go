package Mdown

import (
	"testing"
	"time"
	"sync"
)

func Test_Download(t *testing.T) {
	src := "http://abc.com/abc.zip"
	wg:=&sync.WaitGroup{}  //lock
	timeout:=30*time.Second

	MultiDownload(src, "a.zip", timeout,wg)
}
