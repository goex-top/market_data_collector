package influxdb

import (
	"context"
	"fmt"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/nntaoli-project/goex"
	"log"
	"strings"
	"time"
)

type InfluxdbStorage struct {
	ctx          context.Context
	exchangeName string
	pair         string
	contractType string
	tag          map[string]string
	Url,
	DatabaseName,
	Username,
	Password string
	cli            client.Client
	saveDepthChan  chan goex.Depth
	saveTickerChan chan goex.Ticker
	saveKlineChan  chan goex.Kline
}

func NewInfluxdb(ctx context.Context,
	exchangeName string,
	pair string,
	contractType string,
	url,
	databaseName,
	username,
	password string,
) *InfluxdbStorage {
	s := &InfluxdbStorage{
		ctx:          ctx,
		Username:     username,
		Password:     password,
		Url:          url,
		DatabaseName: databaseName,
		exchangeName: exchangeName,
		pair:         pair,
		contractType: contractType,
	}
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     url,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	s.tag = make(map[string]string)
	if strings.Contains(exchangeName, "Future") {
		s.tag["future_"+contractType] = pair
	} else if strings.Contains(exchangeName, "Swap") {
		s.tag["swap"] = pair
	} else {
		s.tag["spot"] = pair
	}
	s.cli = cli
	s.saveTickerChan = make(chan goex.Ticker)
	s.saveDepthChan = make(chan goex.Depth)
	s.saveKlineChan = make(chan goex.Kline)
	return s
}

func (s *InfluxdbStorage) SaveDepth(depth *goex.Depth) {
	if s.saveDepthChan == nil {
		return
	}
	s.saveDepthChan <- *depth
}

func (s *InfluxdbStorage) SaveTicker(ticker *goex.Ticker) {
	if s.saveTickerChan == nil {
		return
	}

	s.saveTickerChan <- *ticker
}

func (s *InfluxdbStorage) SaveKline(kline *goex.Kline) {
	if s.saveKlineChan == nil {
		return
	}

	s.saveKlineChan <- *kline
}

func (s *InfluxdbStorage) Close() {
	//close(s.saveDepthChan)
	//close(s.saveTickerChan)
	//close(s.saveKlineChan)
	s.cli.Close()
}

func (s *InfluxdbStorage) SaveWorker() {
	/*
		MEASUREMENT    | TAGS         | FIELDS
		exchangeName_ticker  | spot=pair    | xxx
		exchangeName_kline  | future_contractType=pair  | xxx
		exchangeName_depth  | swap=pair    | xxx
	*/

	for {
		select {
		case o := <-s.saveDepthChan:
			fields := make(map[string]interface{})
			fields["ts"] = o.UTime.UnixNano() / int64(time.Millisecond) //unit ms
			for k, v := range o.AskList {
				fields[fmt.Sprintf("ask%d_price", k)] = v.Price
				fields[fmt.Sprintf("ask%d_amount", k)] = v.Amount
			}
			for k, v := range o.BidList {
				fields[fmt.Sprintf("bid%d_price", k)] = v.Price
				fields[fmt.Sprintf("bid%d_amount", k)] = v.Amount
			}
			s.WritesPoints(s.exchangeName+"_"+"depth", s.tag, fields)
		case o := <-s.saveTickerChan:
			fields := make(map[string]interface{})
			fields["ts"] = int64(o.Date)
			fields["last"] = o.Last
			fields["buy"] = o.Buy
			fields["sell"] = o.Sell
			fields["vol"] = o.Vol
			fields["high"] = o.High
			fields["low"] = o.Low

			s.WritesPoints(s.exchangeName+"_"+"ticker", s.tag, fields)

		case o := <-s.saveKlineChan:
			fields := make(map[string]interface{})
			fields["ts"] = o.Timestamp
			fields["open"] = o.Open
			fields["high"] = o.High
			fields["low"] = o.Low
			fields["close"] = o.Close
			fields["vol"] = o.Vol
			s.WritesPoints(s.exchangeName+"_"+"kline", s.tag, fields)

		case <-s.ctx.Done():
			s.Close()
			log.Printf("(%s) %s saveWorker exit\n", s.exchangeName, s.pair)
			return
		}
	}
}

//Insert
func (s *InfluxdbStorage) WritesPoints(table string, tags map[string]string, fields map[string]interface{}) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  s.DatabaseName,
		Precision: "ms",
	})

	if err != nil {
		log.Fatal(err)
	}

	pt, err := client.NewPoint(
		table,
		tags,
		fields,
		time.Now(),
	)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	if err := s.cli.Write(bp); err != nil {
		log.Fatal(err)
	}
}

//query
func (s *InfluxdbStorage) QueryDB(cli client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: s.DatabaseName,
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
