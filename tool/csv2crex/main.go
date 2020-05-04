package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/goex-top/market_data_collector/parser"
	csvsto "github.com/goex-top/market_data_collector/storage/csv"
	"os"
	"strconv"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, `csv2crex version: v1.0.0
Usage: csv2crex [-h] [-p path/xxx] [-f depth_xxx.csv]
Options:
`)
	flag.PrintDefaults()
}

func main() {
	var p string
	var f string
	var o string
	var help bool
	flag.StringVar(&p, "p", "./", "path")
	flag.StringVar(&f, "f", "*.csv", "csv file")
	flag.StringVar(&o, "o", "./", "output path")
	flag.BoolVar(&help, "h", false, "this help")
	flag.Usage = usage

	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if p[len(p)-1] != '/' {
		p += "/"
	}
	fmt.Printf("input path:%s\nfile:%s\noutput path:%s\n", p, f, o)
	var target = f
	if p == o {
		target = "new_" + target
	}
	csvParser := parser.NewDepthDataParser(p)
	depth, err := csvParser.Load(f)
	if err != nil {
		fmt.Printf("load file(%s) fail:%v", p+f, err)
		return
	}
	if o[len(o)-1] != '/' {
		o += "/"
	}
	isNew, targetCsvFile := csvsto.OpenCsvFile(o + target)
	targetCsv := csv.NewWriter(targetCsvFile)
	if isNew {
		data := []string{"t"}
		for k := range depth[0].AskList {
			if k >= 20 {
				break
			}
			data = append(data, fmt.Sprintf("asks[%d].price", k))
			data = append(data, fmt.Sprintf("asks[%d].amount", k))
		}
		for k := range depth[0].BidList {
			if k >= 20 {
				break
			}
			data = append(data, fmt.Sprintf("bids[%d].price", k))
			data = append(data, fmt.Sprintf("bids[%d].amount", k))
		}
		targetCsv.Write(data)
		targetCsv.Flush()
	}

	for k := range depth {
		data := []string{strconv.Itoa(int(depth[k].UTime.UnixNano() / int64(time.Millisecond)))}
		for kk, asks := range depth[k].AskList {
			if kk >= 20 {
				break
			}
			data = append(data, strconv.FormatFloat(asks.Price, 'f', -1, 64))
			data = append(data, strconv.FormatFloat(asks.Amount, 'f', -1, 64))
		}
		for kk, bids := range depth[k].BidList {
			if kk >= 20 {
				break
			}
			data = append(data, strconv.FormatFloat(bids.Price, 'f', -1, 64))
			data = append(data, strconv.FormatFloat(bids.Amount, 'f', -1, 64))
		}
		targetCsv.Write(data)
	}
	targetCsv.Flush()
}
