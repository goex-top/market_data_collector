package csv

import (
	"github.com/goex-top/market_center"
	"testing"
	"time"
)

var s = &CsvStorage{prefix: "../output/csv", outputPath: "../output/tar", exchangeName: market_center.BINANCE, pair: "BTC_USDT"}

func TestCsvStorage_Compress(t *testing.T) {
	ts, _ := time.Parse("2006-01-02", "2020-01-29")

	s.compress(ts)
}
