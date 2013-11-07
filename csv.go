package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"path/filepath"
	//"io"
	"os"
	"runtime"
	"strings"
	"time"
)

const chanLen = 4

var c chan int
var end chan int
var index int

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) //设置cpu的核的数量，从而实现高并发
	c = make(chan int, chanLen)          //简单采用goruntime缓冲，同时最多4个执行，和CPU的数量一致。
	end = make(chan int, 1)
	index = 0

	db, err := dbConnect()

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	path := "./data/163/chddata/sh/20131101"
	filepath.Walk(path, func(path string, f os.FileInfo, e error) error {
		if f == nil {
			return e
		}
		if f.IsDir() {
			return nil
		}
		c <- 0
		go func(db *sql.DB, path string) {
			parseCsv(db, path)
		}(db, path)

		return nil
	})

	<-end
	fmt.Println("scuess!")
}

func dbConnect() (*sql.DB, error) {
	return sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/trackers")
}

func parseCsv(db *sql.DB, path string) {
	index++
	fmt.Println(index)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	ss, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	sz := len(ss)
	for i := 1; i < sz; i++ {
		row := ss[i]
		stmtIns, err := db.Prepare("INSERT INTO 163_transaction_log(stockCode, dateTime, openPrice, highPrice, lowPrice, closePrice, adjClosePrice, volume) VALUES (?, ?, ?, ?, ?, ?, ?, ?)") // ? = placeholder
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

		dateTime, _ := time.Parse("2006-01-02", row[0])
		intTime := dateTime.Unix()
		stockCode := strings.Replace(row[1], "'", "", 0) //还存在编码问题，因为csv文件是ANSI格式的，还要转码才行。
		_, err = stmtIns.Exec(stockCode, intTime, row[6], row[4], row[5], row[3], row[7], row[11])
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	<-c
	if index == 941 {
		end <- 0
	}
}
