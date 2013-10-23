//批量获取雅虎股票数据。
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	UA = "Golang Downloader from Ijibu.com"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //设置cpu的核的数量，从而实现高并发
	c := make(chan bool)
	fh, ferr := os.Open("./shang.ini")
	if ferr != nil {
		return
	}
	defer fh.Close()
	inputread := bufio.NewReader(fh)

	for i := 0; i < 1162; i++ {
		input, _ := inputread.ReadString('\n')
		go getShangTickerTables(c, i, strings.TrimSpace(input))
	}

	<-c

	fmt.Println("main ok")
}

func getShangTickerTables(c chan bool, n int, code string) {
	fileName := "./data/sh/" + code + ".csv"
	f, err := os.OpenFile(fileName, os.O_CREATE, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
	if err != nil {
		panic(err)
	}
	stat, err := f.Stat() //获取文件状态
	if err != nil {
		panic(err)
	}
	defer f.Close()
	urls := "http://table.finance.yahoo.com/table.csv?s=" + code + ".ss"
	var req http.Request
	req.Method = "GET"
	req.Close = true
	req.URL, err = url.Parse(urls)

	if err != nil {
		panic(err)
	}
	header := http.Header{}
	header.Set("Range", "bytes="+strconv.Itoa(int(stat.Size()))+"-")
	header.Set("User-Agent", UA)
	req.Header = header
	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		io.Copy(f, resp.Body)
	}

	if n == 1161 {
		c <- true
	}
	defer resp.Body.Close()
}
