package Mdown

import (
	"testing"
)

func Test_Download(t *testing.T) {
	origin := "http://mirrors.163.com/centos/7.6.1810/isos/x86_64/CentOS-7-x86_64-LiveGNOME-1810.iso"
	target := "a.iso"
	timeout := 30 // second

	downloader := NewMDownloader(origin, target, timeout)
	downloader.Start()
}
