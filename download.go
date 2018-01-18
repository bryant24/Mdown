package Mdown

import (
	"time"
	"strconv"
	"net/http"
	"sync"
	"os"
	"log"
	"io"
)

func MultiDownload(src, target string, timeout time.Duration, wg *sync.WaitGroup) (err error) {

	var res *http.Response
	for {
		res, err = http.Head(src)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	maps := res.Header
	length, err := strconv.Atoi(maps["Content-Length"][0]) // Get the content length from the header request
	if err != nil {
		return
	}

	//根据文件大小分线程数量
	thread := length / 10485760
	if thread == 0 {
		thread = 1
	}
	if thread > 20 && thread < 50 {
		thread = 15
	}
	if thread >= 50 {
		thread = 20
	}

	len_sub := length / thread // Bytes for each Go-routine
	diff := length % thread    // Get the remaining for the last request
	for i := 0; i < thread; i++ {
		wg.Add(1)
		min := len_sub * i       // Min range
		max := len_sub * (i + 1) // Max range
		if (i == thread-1) {
			max += diff
		}
		req, _ := http.NewRequest("GET", src, nil)
		var transport http.RoundTripper = &http.Transport{
			DisableKeepAlives: true,
		}
		client := http.Client{
			Transport: transport,
			Timeout:   time.Duration(timeout * time.Second),
		}
		go func(min int, max int, i int) {
			var temp_file *os.File
			for {
				filename := target + "." + strconv.Itoa(i)
				fileinfo, err := os.Stat(filename)
				filesize := 0
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
						filesize = int(fileinfo.Size())
					}
					if min+filesize == max {
						log.Println("chunk ",i,"download finish")
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
				bytesrange := "bytes=" + strconv.Itoa(min+filesize) + "-" + strconv.Itoa(max-1)
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
			wg.Done()
		}(min, max, i)
	}
	wg.Wait()

	os.Remove(target)
	f, _ := os.OpenFile(target, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 777)
	for j := 0; j < thread; j++ {
		chunkf, _ := os.Open(target + "." + strconv.Itoa(j))
		_, err = io.Copy(f, chunkf)
		if err != nil {
			log.Println("[merge error]")
		}
		chunkf.Close()
		os.Remove(target + "." + strconv.Itoa(j))
	}
	f.Close()
	return
}
