package ws

import "github.com/nntaoli-project/goex"

type WS interface {
	SetCallbacks(
		tickerCallback func(ticker *goex.Ticker),
		depthCallback func(*goex.Depth),
		tradeCallback func(*goex.Trade),
		klineCallback func(*goex.Kline, int))
	SubscribeTicker(pair goex.CurrencyPair) error
	SubscribeDepth(pair goex.CurrencyPair, size int) error
	SubscribeTrade(pair goex.CurrencyPair) error
	SubscribeKline(pair goex.CurrencyPair, period int) error
}

type FutureWS interface {
	SetCallbacks(tickerCallback func(ticker *goex.FutureTicker),
		depthCallback func(*goex.Depth),
		tradeCallback func(*goex.Trade, string),
		klineCallback func(*goex.FutureKline, int))
	SubscribeTicker(pair goex.CurrencyPair, contract string) error
	SubscribeDepth(pair goex.CurrencyPair, contract string, size int) error
	SubscribeTrade(pair goex.CurrencyPair, contract string) error
	SubscribeKline(pair goex.CurrencyPair, contract string, period int) error
}
