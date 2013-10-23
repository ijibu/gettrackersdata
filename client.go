//批量获取雅虎股票数据。
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var j int = 0

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //设置cpu的核的数量，从而实现高并发
	c := make(chan bool)
	fh, ferr := os.Open("./shang.ini")
	if ferr != nil {
		return
	}
	defer fh.Close()
	inputread := bufio.NewReader(fh)

	for i := 0; i < 100; i++ {
		input, _ := inputread.ReadString('\n')
		go getShangTickerTables(c, strings.TrimSpace(input))
	}

	<-c

	fmt.Println("main ok")
}

func getShangTickerTables(c chan bool, code string) {
	fileName := "./data/sh/" + code + ".csv"
	f, err := os.OpenFile(fileName, os.O_CREATE, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
	if err != nil {
		panic(err)
	}

	defer f.Close()
	urls := "http://table.finance.yahoo.com/table.csv?s=" + code + ".ss"

	resp, err := http.Get(urls)
	j++
	if err != nil {
		panic(err)
	}
	io.Copy(f, resp.Body)
	//if resp.StatusCode == 200 {
	//	io.Copy(f, resp.Body)
	//}

	if j == 99 {
		c <- true
	}
	defer resp.Body.Close()
}
