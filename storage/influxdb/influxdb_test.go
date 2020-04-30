package influxdb

import (
	"context"
	"fmt"
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/binance"
	"testing"
	"time"
)

var s = NewInfluxdb(context.Background(), "binance", "btc_usdt", "",
	"http://localhost:8086", "testdb", "admin", "")

func TestNewInfluxdb(t *testing.T) {

}

func storeTicker(ticker *goex.Ticker) {
	fmt.Println(ticker)
	s.WritesPoints("okex_future", map[string]string{"this_week": "btc_usd"}, map[string]interface{}{
		"Pair": ticker.Pair.String(),
		"Last": ticker.Last,
		"Buy":  ticker.Buy,
		"Sell": ticker.Sell,
	})

}

func TestInfluxdbStorage_WritesPoints(t *testing.T) {
	ws := binance.NewBinanceWs()
	ws.ProxyUrl("socks5://127.0.0.1:1080")
	ws.SetCallbacks(storeTicker, nil, nil, nil)
	ws.SubscribeTicker(goex.BTC_USDT)
	time.Sleep(time.Minute)
}
