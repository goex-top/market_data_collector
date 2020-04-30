package storage

import "github.com/nntaoli-project/goex"

type Storage interface {
	SaveDepth(depth *goex.Depth)
	SaveTicker(ticker *goex.Ticker)
	SaveKline(kline *goex.Kline)
	SaveWorker()
	Close()
}
