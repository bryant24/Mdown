package Mdown

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type MDownloader struct {
	Origin  string
	Target  string
	Timeout int
	Wg      *sync.WaitGroup
}

func NewMDownloader(origin, target string, timeout int) *MDownloader {
	newDown := new(MDownloader)
	newDown.Origin = origin
	newDown.Target = target
	newDown.Timeout = timeout
	newDown.Wg = &sync.WaitGroup{}

	return newDown
}

func (md *MDownloader) Start() (err error) {
	var res *http.Response
	for {
		res, err = http.Head(md.Origin)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	maps := res.Header
	length, err := strconv.ParseInt(maps["Content-Length"][0], 10, 64)
	if err != nil {
		return
	}

	thread := md.getThread(length)
	len_sub := length / thread
	diff := length % thread

	for i := int64(0); i < thread; i++ {
		md.Wg.Add(1)
		min := len_sub * i       // Min range
		max := len_sub * (i + 1) // Max range
		if (i == thread-1) {
			max += diff
		}
		req, _ := http.NewRequest("GET", md.Origin, nil)
		var transport http.RoundTripper = &http.Transport{
			DisableKeepAlives: true,
		}
		client := http.Client{
			Transport: transport,
			Timeout:   time.Duration(md.Timeout) * time.Second,
		}
		go func(min int64, max int64, i int64) {
			var temp_file *os.File
			for {
				filename := md.Target + "." + strconv.FormatInt(i, 10)
				fileinfo, err := os.Stat(filename)
				filesize := int64(0)
				if err != nil {
					var e error
					temp_file, e = os.Create(filename)
					if e != nil {
						log.Println("create chunk failed", e)
						continue
					}
				} else {
					log.Println("find before chunk", i)
					temp_file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0777)
					if err != nil {
						log.Println("Open chunk file failed", err)
						os.Remove(filename)
					} else {
						filesize = fileinfo.Size()
					}
					if min+filesize == max {
						log.Println("chunk ", i, "download finish")
						temp_file.Close()
						break
					}
					if min+filesize > max {
						temp_file.Close()
						log.Println("start is bigger than finish")
						os.Remove(filename)
						continue
					}
				}
				bytesrange := "bytes=" + strconv.FormatInt(min+filesize, 10) + "-" + strconv.FormatInt(max-1, 10)
				req.Header.Add("Range", bytesrange)

				resp, e := client.Do(req)
				if e != nil {
					log.Println("[request error],retry:", e)
					temp_file.Close()
					continue
				}
				_, e = io.Copy(temp_file, resp.Body)
				if e != nil {
					log.Println("[copy error],retry", e)
					temp_file.Close()
					resp.Body.Close()
					continue
				}
				temp_file.Close()
				resp.Body.Close()
				break
			}
			md.Wg.Done()
		}(min, max, i)
	}
	md.Wg.Wait()

	os.Remove(md.Target)
	f, _ := os.OpenFile(md.Target, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 777)
	for j := int64(0); j < thread; j++ {
		chunkf, _ := os.Open(md.Target + "." + strconv.FormatInt(j, 10))
		_, err = io.Copy(f, chunkf)
		if err != nil {
			log.Println("[merge error]")
		}
		chunkf.Close()
		os.Remove(md.Target + "." + strconv.FormatInt(j, 10))
	}
	f.Close()
	return
}

func (md *MDownloader) Cancel() {

}

func (md *MDownloader) getThread(filesize int64) (thread int64) {
	// thread caculate
	v := filesize / 10485760
	if v == 0 {
		thread = 1
		return
	}
	if v > 20 && v < 50 {
		thread = 15
		return
	}
	if v >= 50 {
		thread = 20
		return
	}
	return v
}
