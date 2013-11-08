package main

import (
	"fmt"
	"github.com/qiniu/iconv"
	"os"
	"io"
)

func main() {
	cd, err := iconv.Open("utf-8", "latin1")
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return
	}
	defer cd.Close()

	output, _ := os.OpenFile("utf8.csv", os.O_RDWR|os.O_CREATE, 0)
	input, err := os.OpenFile("./data/163/chddata/sh/20131101/600004.csv", os.O_RDWR|os.O_CREATE, 0)
	bufSize := 0
	r := iconv.NewReader(cd, input, bufSize)

	_, err = io.Copy(output, r)
	if err != nil {
		fmt.Println("\nio.Copy failed:", err)
		return
	}
}
