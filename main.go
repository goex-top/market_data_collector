package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/goex-top/market_center"
	mcc "github.com/goex-top/market_center/client"
	"github.com/goex-top/market_data_collector/client"
	"github.com/goex-top/market_data_collector/collector"
	"github.com/goex-top/market_data_collector/config"
	"github.com/goex-top/market_data_collector/storage"
	"github.com/goex-top/market_data_collector/storage/csv"
	"github.com/goex-top/market_data_collector/storage/influxdb"
	"github.com/jinzhu/configor"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	cfg config.Config
)

func usage() {
	fmt.Fprintf(os.Stderr, `market data collector version: v1.0.0
Usage: market_data_collector [-h] [-c config.json]
Options:
`)
	flag.PrintDefaults()
}

func main() {
	var c string
	var help bool
	flag.StringVar(&c, "c", "config.yml", "set configuration `json/yml file`")
	flag.BoolVar(&help, "h", false, "this help")
	flag.Usage = usage

	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	err := configor.Load(&cfg, c)

	if err != nil {
		panic(err)
	}

	if cfg.Store.Csv == cfg.Store.InfluxDB {
		panic("currently only support csv, please check your configure")
	}

	ctx, cancel := context.WithCancel(context.Background())
	for _, v := range cfg.Subs {
		var sto storage.Storage
		if cfg.Store.Csv {
			sto = csv.NewCsvStorage(ctx, v.ExchangeName, v.CurrencyPair, v.ContractType, v.Flag, "output/csv", "output/tar")
		}
		if cfg.Store.InfluxDB {
			sto = influxdb.NewInfluxdb(ctx, v.ExchangeName, v.CurrencyPair, v.ContractType, cfg.Store.InfluxDbCfg.Url, cfg.Store.InfluxDbCfg.Database, cfg.Store.InfluxDbCfg.Username, cfg.Store.InfluxDbCfg.Password)
		}
		go sto.SaveWorker()
		cl := &client.Client{}
		if !cfg.WithMarketCenter {
			cl = client.NewClient(v.ExchangeName, v.CurrencyPair, v.ContractType, nil)
		} else {
			mccc := mcc.NewClient()
			isSpot := market_center.IsFutureExchange(v.ExchangeName)
			if v.Flag&market_center.DataFlag_Depth != 0 {
				if isSpot {
					mccc.SubscribeSpotDepth(v.ExchangeName, v.CurrencyPair, v.Period)
				} else {
					mccc.SubscribeFutureDepth(v.ExchangeName, v.ContractType, v.CurrencyPair, v.Period)
				}
			}
			if v.Flag&market_center.DataFlag_Ticker != 0 {
				if isSpot {
					mccc.SubscribeSpotTicker(v.ExchangeName, v.CurrencyPair, v.Period)
				} else {
					mccc.SubscribeFutureTicker(v.ExchangeName, v.ContractType, v.CurrencyPair, v.Period)
				}
			}
			cl = client.NewClient(v.ExchangeName, v.CurrencyPair, v.ContractType, mccc)
		}

		collector.NewCollector(ctx, cl, v.Period, v.Flag, sto)
	}

	exitSignal := make(chan os.Signal, 1)
	sigs := []os.Signal{os.Interrupt, syscall.SIGILL, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM}
	signal.Notify(exitSignal, sigs...)
	<-exitSignal
	cancel()
	time.Sleep(time.Second)
	log.Println("market data collector exit")
}
