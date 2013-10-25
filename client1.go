//批量获取雅虎股票数据。
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	UA = "Golang Downloader from Ijibu.com"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //设置cpu的核的数量，从而实现高并发
	logfile, _ := os.OpenFile("./test.log", os.O_RDWR|os.O_CREATE, 0)
	c := make(chan bool)
	fh, ferr := os.Open("./shen.ini")
	if ferr != nil {
		return
	}
	defer fh.Close()
	inputread := bufio.NewReader(fh)

	for i := 0; i < 192; i++ { //加入goroutine缓冲，4个执行完了再执行下面的4个
		for k := 0; k < 10; k++ {
			input, _ := inputread.ReadString('\n')
			go getShangTickerTables(logfile, c, i, strings.TrimSpace(input))
		}
		<-c
	}

	for {
		time.Sleep(60000 * time.Second)
		fmt.Println("main ok")
	}
}

func getShangTickerTables(logfile *os.File, c chan bool, n int, code string) {
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	defer logfile.Close()
	flg := make(chan bool, 1)
	fileName := "./data/sz/" + code + ".csv"
	f, err := os.OpenFile(fileName, os.O_CREATE, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
	if err != nil {
		panic(err)
	}

	defer f.Close()

	go func(logfile *os.File, f *os.File, n int, code string) {
		urls := "http://table.finance.yahoo.com/table.csv?s=" + code + ".ss"
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
				logger.Println(code + ":sucess" + strconv.Itoa(resp.StatusCode))
				io.Copy(f, resp.Body)
			} else {
				logger.Println(code + ":http get StatusCode" + strconv.Itoa(resp.StatusCode))
				f.WriteString("http get StatusCode" + strconv.Itoa(resp.StatusCode))
			}
			defer resp.Body.Close()
		} else {
			logger.Println(code + ":http get error" + code)
			f.WriteString("http get error")
		}
		flg <- true
	}(logfile, f, n, code)

	//加入超时设置
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(500 * time.Second)
		timeout <- true
	}()

	select {
	case <-flg:
		//正常退出
		if n%10 == 0 {
			c <- true
		}
		return
	case <-timeout:
		//超时退出
		if n%10 == 0 {
			c <- true
		}
		logger.Println(code + ":http get timeout" + code)
		f.WriteString("timeouttimeouttimeouttimeouttimeouttimeouttimeouttimeouttimeouttimeouttimeouttimeout")
		return
	}

}
