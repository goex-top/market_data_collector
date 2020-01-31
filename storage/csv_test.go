package storage

import (
	"testing"
	"time"
)

var s = &CsvStorage{prefix: "../csv", outputPath: "../tar", exchangeName: "binance.com", pair: "BTC_USDT"}

func TestCsvStorage_Compress(t *testing.T) {
	ts, _ := time.Parse("2006-01-02", "2020-01-29")

	s.compress(ts)
}
