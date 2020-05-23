package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	csvsto "github.com/goex-top/market_data_collector/storage/csv"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, `influx2csv version: v1.0.0
Usage: influx2csv [-h] [-p path/xxx] [-f depth_xxx.csv]
Options:
`)
	flag.PrintDefaults()
}

//Url,
//DatabaseName,
//Username,
//Password string
//exchangeName string
//pair         string
//contractType string

func main() {
	var databasename string
	var username string
	var password string
	var url string
	var tag string
	var o string
	var start string
	var end string
	var help bool
	flag.StringVar(&o, "o", "./", "output path")
	flag.StringVar(&url, "u", "", "influxdb url")
	flag.StringVar(&databasename, "d", "", "database name")
	flag.StringVar(&username, "n", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.StringVar(&tag, "t", "", "tag, Future_Okex_next_week=BTC_USDT")
	flag.StringVar(&start, "s", "", "start time, 2020-05-22T23:00:00Z")
	flag.StringVar(&end, "e", "", "end time, 2020-05-22T23:00:00Z")
	flag.BoolVar(&help, "h", false, "this help")
	flag.Usage = usage

	flag.Parse()

	//url = "http://localhost:8086"
	databasename = "market_data"
	if help {
		flag.Usage()
		return
	}
	if url == "" {
		fmt.Println("please input url")
		return
	}
	if databasename == "" {
		fmt.Println("please input databasename")
		return
	}
	if tag == "" {
		fmt.Println("please input tag")
		return
	}

	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     url,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}

	// verify tag
	t := strings.Split(tag, "=")
	q := fmt.Sprintf("SHOW TAG VALUES ON %s FROM depth WITH KEY = %s", databasename, t[0])
	ret, err := QueryDB(cli, databasename, q)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(ret) == 0 || len(ret[0].Series) == 0 {
		fmt.Println("tag query no data")
		return
	}

	//start = "2020-05-22T23:00:00Z"
	//end = "2020-05-23T23:05:00Z"

	// start query data

	q = "SELECT last(ts) AS ts,"
	for i := 0; i < 20; i++ {
		q += fmt.Sprintf(" last(\"ask%d_price\") AS \"asks[%d].price\", last(\"ask%d_amount\") AS \"asks[%d].amount\",", i, i, i, i)
	}
	for i := 0; i < 20; i++ {
		q += fmt.Sprintf(" last(\"bid%d_price\") AS \"bids[%d].price\", last(\"bid%d_amount\") AS \"bids[%d].amount\",", i, i, i, i)
	}
	q = q[:len(q)-1]
	q += fmt.Sprintf(" FROM \"%s\".\"autogen\".\"depth\" WHERE \"%s\"='%s' AND time >= '%s' AND time <= '%s' GROUP BY time(1s) FILL(null)", databasename, t[0], t[1], start, end)

	fmt.Println(q)
	ret, err = QueryDB(cli, databasename, q)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(ret) == 0 || len(ret[0].Series) == 0 {
		fmt.Println("query no data ", ret[0].Err)
		return
	}
	if o[len(o)-1] != '/' {
		o += "/"
	}
	target := fmt.Sprintf("depth_%s_%s_%s_%s.csv", t[0], t[1], start, end)
	isNew, targetCsvFile := csvsto.OpenCsvFile(o + target)
	targetCsv := csv.NewWriter(targetCsvFile)
	if isNew {
		data := []string{"t"}
		for k := 0; k < 20; k++ {
			data = append(data, fmt.Sprintf("asks[%d].price", k))
			data = append(data, fmt.Sprintf("asks[%d].amount", k))
		}
		for k := 0; k < 20; k++ {
			data = append(data, fmt.Sprintf("bids[%d].price", k))
			data = append(data, fmt.Sprintf("bids[%d].amount", k))
		}
		targetCsv.Write(data)
		targetCsv.Flush()
	}

	for _, row := range ret[0].Series {
		for _, value := range row.Values {
			if value[1] != nil {
				data := make([]string, 0)
				data = append(data, string(value[1].(json.Number)))
				//value = value[2:]
				for k := 2; k < len(value); k += 2 {
					data = append(data, string(value[k].(json.Number)))
					data = append(data, string(value[k+1].(json.Number)))
				}
				targetCsv.Write(data)
			}
		}
	}
	targetCsv.Flush()
}

//query
func QueryDB(cli client.Client, database, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: database,
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
