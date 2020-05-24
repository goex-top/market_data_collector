package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	csvsto "github.com/goex-top/market_data_collector/storage/csv"
	client "github.com/influxdata/influxdb1-client/v2"
	"os"
	"strings"
	"time"
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
	flag.StringVar(&url, "u", "http://localhost:8086", "influxdb url, http://localhost:8086")
	flag.StringVar(&databasename, "d", "market_data", "database name, market_data")
	flag.StringVar(&username, "n", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.StringVar(&tag, "t", "", "tag, Future_Okex_next_week=BTC_USDT")
	flag.StringVar(&start, "s", "", "start time, 2020-05-22T23:00:00Z")
	flag.StringVar(&end, "e", "", "end time, 2020-05-22T23:00:00Z")
	flag.BoolVar(&help, "h", false, "this help")
	flag.Usage = usage

	flag.Parse()

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
	fmt.Printf("url:%s\ndatabase:%s\ntag:%s\nstart:%s\nend:%s\n", url, databasename, tag, start, end)
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     url,
		Username: username,
		Password: password,
	})
	if err != nil {
		fmt.Printf("new influxdb client fail:%s, exit!\n", err)
		return
	}
	_, _, err = cli.Ping(time.Second * 5)
	if err != nil {
		fmt.Printf("ping infuxdb fail:%s, exit!\n", err)
		return
	}
	fmt.Printf("ping %s is ok\n", url)
	if o[len(o)-1] != '/' {
		o += "/"
	}
	t := strings.Split(tag, "=")
	target := fmt.Sprintf("depth_%s_%s_%s_%s.csv", t[0], t[1], start, end)
	isNew, targetCsvFile := csvsto.OpenCsvFile(o + target)
	defer targetCsvFile.Close()

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
	} else {
		fmt.Printf("file exist:%s, exit!\n", target)
		return
	}

	// verify tag
	q := fmt.Sprintf("SHOW TAG VALUES ON %s FROM depth WITH KEY = %s", databasename, t[0])
	fmt.Println(q)
	ret, err := QueryDB(cli, databasename, q)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(ret) == 0 || len(ret[0].Series) == 0 {
		fmt.Println("tag query no data")
		return
	}

	// start query data
	st, err := time.Parse("2006-01-02T15:04:05Z", start)
	if err != nil {
		fmt.Printf("start time(%s) format error:%s\n", start, err)
		return
	}
	et, err := time.Parse("2006-01-02T15:04:05Z", end)
	if err != nil {
		fmt.Printf("end time(%s) format error:%s\n", end, err)
		return
	}
	minutes := 10
	count := 0
	n := int(et.Sub(st).Minutes()/float64(minutes)) + 1
	for k := 0; k < n; k++ {
		q = "SELECT last(ts) AS ts,"
		for i := 0; i < 20; i++ {
			q += fmt.Sprintf(" last(\"ask%d_price\") AS \"asks[%d].price\", last(\"ask%d_amount\") AS \"asks[%d].amount\",", i, i, i, i)
		}
		for i := 0; i < 20; i++ {
			q += fmt.Sprintf(" last(\"bid%d_price\") AS \"bids[%d].price\", last(\"bid%d_amount\") AS \"bids[%d].amount\",", i, i, i, i)
		}
		q = q[:len(q)-1]
		q += fmt.Sprintf(" FROM \"%s\".\"autogen\".\"depth\" WHERE \"%s\"='%s' AND time >= '%s' AND time <= '%s' GROUP BY time(1ms) FILL(null)",
			databasename, t[0], t[1], st.Format("2006-01-02T15:04:05Z"), st.Add(time.Hour*4).Format("2006-01-02T15:04:05Z"))

		st = st.Add(time.Duration(minutes) * time.Minute)

		fmt.Println(q)
		ret, err = QueryDB(cli, databasename, q)
		if err != nil {
			fmt.Printf("[%s - %s] error:%s\n", st.Format("2006-01-02T15:04:05Z"), st.Add(time.Hour*4).Format("2006-01-02T15:04:05Z"), err)
			continue
		}
		if len(ret) == 0 || len(ret[0].Series) == 0 {
			fmt.Printf("[%s - %s] query no data %s\n", st.Format("2006-01-02T15:04:05Z"), st.Add(time.Hour*4).Format("2006-01-02T15:04:05Z"), ret[0].Err)
			continue
		}

		for _, row := range ret[0].Series {
			for _, value := range row.Values {
				if value[1] != nil {
					count++
					data := make([]string, 0)
					data = append(data, string(value[1].(json.Number)))
					for k := 2; k < len(value); k += 2 {
						data = append(data, string(value[k].(json.Number)))
						data = append(data, string(value[k+1].(json.Number)))
					}
					targetCsv.Write(data)
				}
			}
		}
		targetCsv.Flush()
		fmt.Printf("%d rows has been exported sucuessful!\n", len(ret[0].Series))
	}

	fmt.Printf("\ntotal %d rows has been exported to %s sucuessful!\n", count, target)
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
