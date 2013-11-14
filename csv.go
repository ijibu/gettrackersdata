package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var num *int = flag.Int("n", 0, "please input a num like 1024")
var stockType *string = flag.String("s", "sh", "please input a stockType like sh")
var dataType *string = flag.String("t", "chddata", "please input a dataType like chddata")

const chanLen = 4

var c chan int
var end chan int
var index int
var path string

func main() {
	flag.Parse()

	if *num == 0 {
		show_usage()
		return
	}

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

	if *dataType == "cjmx" {
		path = "./data/163/chddata/" + *stockType + "/20131101"
	} else if *dataType == "chddata" {
		path = "./data/163/chddata/" + *stockType + "/20131101"
	} else if *dataType == "lszjlx" {
		path = "./data/163/lszjlx/" + *stockType + "/csv"
	}

	filepath.Walk(path, func(path string, f os.FileInfo, e error) error {
		if f == nil {
			return e
		}
		if f.IsDir() {
			return nil
		}
		c <- 0

		if *dataType == "cjmx" {
			go parseCsv(db, path)
		} else if *dataType == "chddata" {
			go parseCsv(db, path)
		} else if *dataType == "lszjlx" {
			go importLszjlxCsv(db, path)
		}

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
		stockCode := row[1] //直接处理CSV文件，不用在程序里面进行替换了
		_, err = stmtIns.Exec(stockCode, intTime, row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9])
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	<-c
	if index == 941 {
		end <- 0
	}
}

/**
 * 导入历史资金流向数据
 */
func importLszjlxCsv(db *sql.DB, path string) {
	index++
	fmt.Println(index)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	stockCode := getFileName(path)

	reader := csv.NewReader(file)
	ss, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	sz := len(ss)
	for i := 1; i < sz; i++ {
		row := ss[i]
		stmtIns, err := db.Prepare("INSERT INTO lszjlx(stockCode, dateTime, closePrice, changePrice, percentChg, income, expenditure, netIncome, mainIncome, mainExpenditure, netMainIncome) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

		dateTime, _ := time.Parse("2006-01-02", row[0])
		intTime := dateTime.Unix()
		_, err = stmtIns.Exec(stockCode, intTime, row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9])
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	<-c
	if index == *num {
		end <- 0
	}
}

func show_usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [-d=<date>] [-n=<num>] [-s=<stockType>] [-t=<type>]\n"+
			"       <command> [<args>]\n\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Flags:\n")
	flag.PrintDefaults()
}

/**
 * 根据路径名获取文件名
 */
func getFileName(fileFullPath string) string {
	fName := filepath.Base(fileFullPath)
	extName := filepath.Ext(fileFullPath)
	bName := fName[:len(fName)-len(extName)]
	return bName
}
