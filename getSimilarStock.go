//利用余弦相似性计算股票的K线相似性
//package CosineSimilarity
/*
	使用说明：
	go build getSimilarStock.go
	./getSimilarStock -n=[股票代码] > 文件名(输出结果到那个文件)
*/

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

var stockCode *int = flag.Int("n", 0, "please input a stockCode like 600000")

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

type rateArr []float64

func main() {
	var (
		rate    map[int]rateArr = make(map[int]rateArr, 943)
		rate1   []float64
		stock   int
		path    string = "./data/163/chddata/sh/20131101"
		cosine  float64
		retData ByCosine
	)

	flag.Parse()

	if *stockCode != 0 { //获取某只股票的相似性
		fileName := "./data/163/chddata/sh/20131101/" + strconv.Itoa(*stockCode) + ".csv"
		_, rate2 := getRateFromCsv(fileName)
		rate[*stockCode] = rate2
	}

	filepath.Walk(path, func(path string, f os.FileInfo, e error) error {
		if f == nil {
			return e
		}
		if f.IsDir() {
			return nil
		}
		stock, rate1 = getRateFromCsv(path)
		rate[stock] = rate1
		return nil
	})

	//此算法还有一个bug，比如计算股票A和股票B的相似性计算了两次。实际上该处需要的是一种组合，而不是排列。
	//所以采用组合算法可以满足需求。
	for code, rateItem := range rate {
		if *stockCode != 0 {
			if code != *stockCode {
				cosine = CosDist(rate[*stockCode], rateItem)
				retData = append(retData, stockDist{*stockCode, code, cosine})
			}
		} else {
			for code1, rateItem1 := range rate {
				if code1 != code {
					cosine = CosDist(rateItem, rateItem1)
					retData = append(retData, stockDist{code, code1, cosine})
				}
			}
		}
	}

	sort.Sort(ByCosine(retData))
	for i := 0; i < len(retData); i++ {
		if i%2 == 0 || retData[i].cosine == 0 {
			continue
		}
		fmt.Println(retData[i])
	}
}

//从csv文件中读取数据，组成向量
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
		"Usage: %s  [-n=<num>]\n"+
			"       <command> [<args>]\n\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Flags:\n")
	flag.PrintDefaults()
}
