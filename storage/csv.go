package storage

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/goex-top/market_center"
	jsoniter "github.com/json-iterator/go"
	"github.com/nntaoli-project/goex"
	"log"
	"os"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type CsvStorage struct {
	exchangeName   string
	pair           string
	contractType   string
	flag           market_center.DataFlag
	prefix         string
	outputPath     string
	saveDepthChan  chan goex.Depth
	saveTickerChan chan goex.Ticker
	saveKlineChan  chan goex.Kline
	fileTimestamp  time.Time
	ctx            context.Context
	depthFile      *os.File
	tickerFile     *os.File
	klineFile      *os.File
	depthCsv       *csv.Writer
	tickerCsv      *csv.Writer
	klineCsv       *csv.Writer
}

func NewCsvStorage(
	ctx context.Context,
	exchangeName string,
	pair string,
	contractType string,
	flag market_center.DataFlag,
	prefix string,
	outputPath string,
) *CsvStorage {
	var saveDepthChan chan goex.Depth
	var saveTickerChan chan goex.Ticker
	var saveKlineChan chan goex.Kline
	var depthFile *os.File
	var tickerFile *os.File
	var klineFile *os.File
	var depthCsv *csv.Writer
	var tickerCsv *csv.Writer
	var klineCsv *csv.Writer

	fileTimestamp := time.Now()
	ts := fileTimestamp.Format("2006-01-02")
	isNew := false

	if flag&market_center.DataFlag_Depth != 0 {
		isNew, depthFile = openFile(fmt.Sprintf("%s/depth_%s_%s_%s%s.csv", prefix, exchangeName, pair, contractType, ts))
		depthCsv = csv.NewWriter(depthFile)
		if isNew {
			data := []string{"t", "a", "b"}
			depthCsv.Write(data)
			depthCsv.Flush()
		}
		saveDepthChan = make(chan goex.Depth)
	}

	if flag&market_center.DataFlag_Ticker != 0 {
		isNew, tickerFile = openFile(fmt.Sprintf("%s/ticker_%s_%s_%s%s.csv", prefix, exchangeName, pair, contractType, ts))
		tickerCsv = csv.NewWriter(tickerFile)
		if isNew {
			data := []string{"t", "b", "s", "h", "l", "v"}
			tickerCsv.Write(data)
			tickerCsv.Flush()
		}
		saveTickerChan = make(chan goex.Ticker)
	}

	if flag&market_center.DataFlag_Kline != 0 {
		isNew, klineFile = openFile(fmt.Sprintf("%s/depth_%s_%s_%s%s.csv", prefix, exchangeName, pair, contractType, ts))
		klineCsv = csv.NewWriter(klineFile)
		if isNew {
			data := []string{"t", "o", "h", "l", "c", "v"}
			klineCsv.Write(data)
			klineCsv.Flush()
		}
		saveKlineChan = make(chan goex.Kline)
	}

	return &CsvStorage{
		ctx:            ctx,
		exchangeName:   exchangeName,
		pair:           pair,
		contractType:   contractType,
		flag:           flag,
		prefix:         prefix,
		outputPath:     outputPath,
		saveDepthChan:  saveDepthChan,
		saveTickerChan: saveTickerChan,
		saveKlineChan:  saveKlineChan,
		fileTimestamp:  fileTimestamp,
		depthFile:      depthFile,
		tickerFile:     tickerFile,
		klineFile:      klineFile,
		depthCsv:       depthCsv,
		tickerCsv:      tickerCsv,
		klineCsv:       klineCsv,
	}
}

func (s *CsvStorage) SaveDepth(depth *goex.Depth) {
	if s.saveDepthChan == nil {
		return
	}
	s.saveDepthChan <- *depth
}

func (s *CsvStorage) SaveTicker(ticker *goex.Ticker) {
	if s.saveTickerChan == nil {
		return
	}

	s.saveTickerChan <- *ticker
}

func (s *CsvStorage) SaveKline(kline goex.Kline) {
	if s.saveKlineChan == nil {
		return
	}

	s.saveKlineChan <- kline
}

func openFile(fileName string) (bool, *os.File) {
	var file *os.File
	var err1 error
	var isNew = false
	checkFileIsExist := func(fileName string) bool {
		var exist = true
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			exist = false
		}
		return exist
	}
	if checkFileIsExist(fileName) {
		file, err1 = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 666)
	} else {
		file, err1 = os.Create(fileName)
		isNew = true
	}
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook %v", err1)
		panic(err1)
	}
	return isNew, file
}

func (s *CsvStorage) Close() {
	if s.depthCsv != nil {
		s.depthCsv.Flush()
		s.depthFile.Close()
	}
	if s.tickerCsv != nil {
		s.tickerCsv.Flush()
		s.tickerFile.Close()
	}
	if s.klineCsv != nil {
		s.klineCsv.Flush()
		s.klineFile.Close()
	}
}

func (s *CsvStorage) compress(fileTimestamp time.Time) {
	ts := fileTimestamp.Format("2006-01-02")
	//src := fmt.Sprintf("%s_%s_%s.csv", s.exchangeName, s.pair, ts)
	filters := []string{s.exchangeName, s.pair, ts, ".csv"}
	dst := fmt.Sprintf("%s/%s_%s_%s.tar.gz", s.outputPath, s.exchangeName, s.pair, ts)

	csvs := GetSrcFileName(s.prefix, filters)
	log.Println("start to compress *.csv to *.tar.gz, ts:", ts)
	err := CompressFile(s.prefix, csvs, dst)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range csvs {
		err := os.Remove(s.prefix + "/" + v)
		if err != nil {
			log.Printf("remove file %s fail:%s\n", s.prefix+"/"+v, err.Error())
		} else {
			log.Printf("remove file %s success\n", s.prefix+"/"+v)
		}
	}
}

func (s *CsvStorage) reNewFile() {
	now := time.Now()
	if now.Day() == s.fileTimestamp.Day() {
		return
	}
	s.Close()
	log.Printf("now day:%d, file timestamp day:%d", now.Day(), s.fileTimestamp.Day())
	go s.compress(s.fileTimestamp)

	s.fileTimestamp = now

	ts := s.fileTimestamp.Format("2006-01-02")
	isNew := false

	if s.flag&market_center.DataFlag_Depth != 0 {
		isNew, s.depthFile = openFile(fmt.Sprintf("%s/depth_%s_%s_%s%s.csv", s.prefix, s.exchangeName, s.pair, s.contractType, ts))
		s.depthCsv = csv.NewWriter(s.depthFile)
		if isNew {
			data := []string{"t", "a", "b"}
			s.depthCsv.Write(data)
			s.depthCsv.Flush()
		}
	}

	if s.flag&market_center.DataFlag_Ticker != 0 {
		isNew, s.tickerFile = openFile(fmt.Sprintf("%s/ticker_%s_%s_%s%s.csv", s.prefix, s.exchangeName, s.pair, s.contractType, ts))
		s.tickerCsv = csv.NewWriter(s.tickerFile)
		if isNew {
			data := []string{"t", "b", "s", "h", "l", "v"}
			s.tickerCsv.Write(data)
			s.tickerCsv.Flush()
		}
	}

	if s.flag&market_center.DataFlag_Kline != 0 {
		isNew, s.klineFile = openFile(fmt.Sprintf("%s/depth_%s_%s_%s%s.csv", s.prefix, s.exchangeName, s.pair, s.contractType, ts))
		s.klineCsv = csv.NewWriter(s.klineFile)
		if isNew {
			data := []string{"t", "o", "h", "l", "c", "v"}
			s.klineCsv.Write(data)
			s.klineCsv.Flush()
		}
	}
}

func (s *CsvStorage) SaveWorker() {

	tick := time.NewTicker(time.Second)

	for {
		select {
		case <-tick.C:
			s.reNewFile()
		case o := <-s.saveDepthChan:
			asks := make([][]float64, 0)
			bids := make([][]float64, 0)
			for _, v := range o.AskList {
				ask := make([]float64, 0)
				ask = append(ask, v.Price, v.Amount)
				asks = append(asks, ask)
			}
			a, _ := json.Marshal(asks)

			for _, v := range o.BidList {
				bid := make([]float64, 0)
				bid = append(bid, v.Price, v.Amount)
				bids = append(bids, bid)
			}
			b, _ := json.Marshal(bids)

			data := []string{
				fmt.Sprint(o.UTime.UnixNano() / int64(time.Millisecond)),
				string(a),
				string(b),
			}

			s.depthCsv.Write(data)
			s.depthCsv.Flush()

		case o := <-s.saveTickerChan:
			data := []string{
				fmt.Sprint(o.Date),
				fmt.Sprint(o.Buy),
				fmt.Sprint(o.Sell),
				fmt.Sprint(o.High),
				fmt.Sprint(o.Low),
				fmt.Sprint(o.Vol),
			}
			s.tickerCsv.Write(data)
			s.tickerCsv.Flush()

		case o := <-s.saveKlineChan:
			data := []string{
				fmt.Sprint(o.Timestamp),
				fmt.Sprint(o.Open),
				fmt.Sprint(o.High),
				fmt.Sprint(o.Low),
				fmt.Sprint(o.Close),
				fmt.Sprint(o.Vol),
			}

			s.klineCsv.Write(data)
			s.klineCsv.Flush()

		case <-s.ctx.Done():
			s.Close()
			log.Printf("(%s) %s saveWorker exit\n", s.exchangeName, s.pair)
			return
		}
	}
}
