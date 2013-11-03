//批量获取雅虎股票数据。
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	UA = "yahoo.com"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //设置cpu的核的数量，从而实现高并发
	logfile, _ := os.OpenFile("./test.log", os.O_RDWR|os.O_CREATE, 0)
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	path := "./data/erro/163"
	c := make(chan int, 145)
	i := 1

	filepath.Walk(path, func(path string, f os.FileInfo, e error) error {
		if f == nil {
			return e
		}
		if f.IsDir() {
			return nil
		}
		str := strings.Split(path, "\\")
		input := strings.Replace(str[3], ".csv", "", -1)
		//fmt.Println(input)
		go func(logger *log.Logger, logfile *os.File, input string) {
			getShangTickerTables(logger, logfile, input)
			c <- 0
		}(logger, logfile, strings.TrimSpace(input))
		//time.Sleep(1 * time.Second) //加入执行缓冲，否则同时发起大量的tcp连接，操作系统会直接返回错误。
		if i%5 == 0 {
			time.Sleep(5 * time.Second) //加入执行缓冲，否则同时发起大量的tcp连接，操作系统会直接返回错误。
		}
		i++
		return nil
	})

	defer logfile.Close()
	for j := 0; j < 145; j++ {
		<-c
	}
	time.Sleep(100 * time.Second) //加入执行缓冲，否则同时发起大量的tcp连接，操作系统会直接返回错误。
}

func getShangTickerTables(logger *log.Logger, logfile *os.File, code string) {
	//并发写文件必须要有锁啊，怎么还是串行程序的思维啊。
	fileName := "./data/163/" + code + ".csv"
	f, err := os.OpenFile(fileName, os.O_CREATE, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//urls := "http://table.finance.yahoo.com/table.csv?s=" + code + ".sz"
	urls := "http://quotes.money.163.com/service/chddata.html?code=0" + code + "&start=19900530&end=20131103&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP"
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
			logger.Println(logfile, code+":sucess"+strconv.Itoa(resp.StatusCode))
			fmt.Println(code + ":sucess" + strconv.Itoa(resp.StatusCode))
			io.Copy(f, resp.Body)
		} else {
			logger.Println(logfile, code+":http get StatusCode"+strconv.Itoa(resp.StatusCode))
			fmt.Println(code + ":http get StatusCode" + strconv.Itoa(resp.StatusCode))
			//删除执行404错误的文件
			//此处会执行失败，为何？因为上面的defer还引用用改文件,或者该文件还是打开的，所以不能删除，那该怎么办呢？
			f.Close()
			//os.Remove(fileName)
			os.Rename(fileName, fileName+"_"+strconv.Itoa(resp.StatusCode))
		}
		defer resp.Body.Close()
	} else {
		logger.Println(logfile, code+":http get error"+code)
		fmt.Println(code + ":http get error" + code)
	}
}
