package main

import (
	"database/sql"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	UA = "Golang Downloader from Ijibu.com"
)

func main() {
	logfile, _ := os.OpenFile("./log/mysql.log", os.O_RDWR|os.O_CREATE, 0)
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/trackers")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	rows, err := db.Query("select exchange, code from stock where status = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		src      string
		urls     string
		exchange int
		code     string
		time     int
	)

	for rows.Next() {
		err := rows.Scan(&exchange, &code)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(exchange, code)
		time = 1383062400
		if exchange == 1 {
			urls = "http://ichart.yahoo.com/table.csv?s=" + code + ".SS&a=09&b=30&c=2013&d=09&e=30&f=2013&g=d"
		} else if exchange == 2 {
			urls = "http://ichart.yahoo.com/table.csv?s=" + code + ".SZ&a=09&b=30&c=2013&d=09&e=30&f=2013&g=d"
		} else {
			urls = "http://ichart.yahoo.com/table.csv?s=" + code + "&a=09&b=29&c=2013&d=09&e=29&f=2013&g=d"
			time = 1382976000
		}
		//log.Println(urls)
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
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println("http read error")
				}
				src = string(body)
				datas := strings.Split(src, "\n")
				row := strings.Split(datas[1], ",")
				if len(row) == 7 {
					stmtIns, err := db.Prepare("INSERT INTO transaction_log(stockCode, dateTime, openPrice, highPrice, lowPrice, closePrice, adjClosePrice, volume) VALUES (?, ?, ?, ?, ?, ?, ?, ?)") // ? = placeholder
					if err != nil {
						panic(err.Error()) // proper error handling instead of panic in your app
					}
					defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

					_, err = stmtIns.Exec(code, time, row[1], row[2], row[3], row[4], row[6], row[5])
					if err != nil {
						panic(err.Error()) // proper error handling instead of panic in your app
					}
					log.Println(code + ":200")
					logger.Println(logfile, code+":sucess200")
				}
			} else {
				logger.Println(logfile, code+":http get StatusCode"+strconv.Itoa(resp.StatusCode))
				log.Println(code + ":" + strconv.Itoa(resp.StatusCode))
			}
			defer resp.Body.Close()
		} else {
			logger.Println(logfile, code+":http get error"+code)
			log.Println(code + ":error")
		}
	}

}
