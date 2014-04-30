//利用余弦相似性计算股票的K线相似性
//package CosineSimilarity
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

//保证对比的向量必须有最小的维度，否则相似性计算出来不准确
const minRateLen int = 2500

var stockCode *int = flag.Int("n", 600633, "please input a stockCode like 600000")

type stockDist struct {
	code1, code2 int
	cosine       float64
}

func (p stockDist) String() string {
	return fmt.Sprintf("%d,%d,%g", p.code1, p.code2, p.cosine)
}

// ByCosine implements sort.Interface for []stockDist based on
// the cosine field.
type ByCosine []stockDist

func (a ByCosine) Len() int           { return len(a) }
func (a ByCosine) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCosine) Less(i, j int) bool { return a[i].cosine > a[j].cosine }

func CosDist(rate1 []float64, rate2 []float64) float64 {
	//if len(rate1) != len(rate2) {
	//	return -0.00001
	//}
	var ratelen int
	if len(rate1) >= len(rate2) {
		ratelen = len(rate2)
	} else {
		ratelen = len(rate1)
	}

	if ratelen < minRateLen {
		return 0
	}

	var (
		sum_xy float64 = 0.0
		sum_x  float64 = 0.0
		sum_y  float64 = 0.0
	)

	for i := 0; i < ratelen; i++ {
		sum_xy += rate1[i] * rate2[i]
		sum_x += rate1[i] * rate1[i]
		sum_y += rate2[i] * rate2[i]
	}

	return sum_xy / math.Sqrt(sum_x*sum_y)
}

func main() {
	var (
		rate1, rate2 []float64
		stock        int
		path         string = "./data/163/chddata/sh/20131101"
		cosine       float64
		retData      ByCosine
	)

	flag.Parse()

	if *stockCode == 0 {
		show_usage()
		return
	}

	fileName := "./data/163/chddata/sh/20131101/" + strconv.Itoa(*stockCode) + ".csv"
	_, rate1 = getRateFromCsv(fileName)

	filepath.Walk(path, func(path string, f os.FileInfo, e error) error {
		if f == nil {
			return e
		}
		if f.IsDir() {
			return nil
		}
		stock, rate2 = getRateFromCsv(path)
		cosine = CosDist(rate1, rate2)
		//fmt.Print(*stockCode)
		//fmt.Print(":")
		//fmt.Print(stock)
		//fmt.Print(": ")
		//fmt.Println(cosine)
		retData = append(retData, stockDist{*stockCode, stock, cosine})

		return nil
	})

	sort.Sort(ByCosine(retData))
	for i := 0; i < len(retData); i++ {
		fmt.Println(retData[i])
	}
}

func getRateFromCsv(path string) (stockCode int, stockRate []float64) {
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
		stockCode, _ = strconv.Atoi(row[1][1:]) //直接处理CSV文件，不用在程序里面进行替换了
		price, _ := strconv.ParseFloat(row[3], 64)
		stockRate = append(stockRate, price)
	}

	return
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
