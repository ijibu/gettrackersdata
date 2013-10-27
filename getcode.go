//批量获取股票交易码。
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	UA = "Golang Downloader from Ijibu.com"
)

func main() {
	getCodes()
}

func getCodes() {
	//并发写文件必须要有锁啊，怎么还是串行程序的思维啊。
	fileName := "./data/data.html"

	f, err := os.Create(fileName) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
	if err != nil {
		panic(err)
	}

	defer f.Close()

	urls := "http://quote.eastmoney.com/stocklist.html"
	var req http.Request
	req.Method = "GET"
	req.Close = true
	req.URL, err = url.Parse(urls)
	if err != nil {
		panic(err)
	}

	header := http.Header{}
	header.Set("User-Agent", UA)
	req.Header = header
	resp, err := http.DefaultClient.Do(&req)
	if err == nil {
		if resp.StatusCode == 200 {
			io.Copy(f, resp.Body)
		} else {
			fmt.Println("http get StatusCode")
		}
		defer resp.Body.Close()
	} else {
		fmt.Println("http get error")
	}
}

func parseHtml() {
	filename := "./data/shanghai.html"
	file, err := os.OpenFile(filename)
	if err != nil {
		fmt.Println("failed to open:", filename)
	}

	defer file.Close()

	finfo, err := file.Stat()
	if err != nil {
		fmt.Println("get file info failed:", file, size)
	}

	size := finfo.Size()
}
