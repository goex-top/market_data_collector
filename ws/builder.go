package ws

import (
	"github.com/goex-top/market_center"
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/binance"
	//"github.com/nntaoli-project/goex/huobi"
	"github.com/nntaoli-project/goex/okex"
)

func BuilderSpot(exchangeName string) WS {
	switch exchangeName {
	case market_center.OKEX:
		return NewOKExSpotV3Ws()
	case market_center.BINANCE:
		return binance.NewBinanceWs()
	default:
		return nil
	}

}

func BuilderFuture(exchangeName string) FutureWS {
	switch exchangeName {
	case market_center.FUTURE_OKEX:
		return NewOKExFutureV3Ws()
	//case market_center.FUTURE_HBDM:
	//	return huobi.NewHbdmWs()
	default:
		return nil
	}
}

type OKExSpotV3Ws struct {
	ok *okex.OKExV3SpotWs
}

func NewOKExSpotV3Ws() *OKExSpotV3Ws {
	return &OKExSpotV3Ws{
		ok: okex.NewOKExSpotV3Ws(nil),
	}
}

func (o *OKExSpotV3Ws) SubscribeTicker(pair goex.CurrencyPair) error {
	return o.ok.SubscribeTicker(pair)
}

func (o *OKExSpotV3Ws) SubscribeDepth(pair goex.CurrencyPair, size int) error {
	return o.ok.SubscribeDepth(pair, size)
}

func (o *OKExSpotV3Ws) SubscribeTrade(pair goex.CurrencyPair) error {
	return o.ok.SubscribeTrade(pair)
}

func (o *OKExSpotV3Ws) SubscribeKline(pair goex.CurrencyPair, period int) error {
	return o.ok.SubscribeKline(pair, period)
}

func (o *OKExSpotV3Ws) SetCallbacks(
	tickerCallback func(ticker *goex.Ticker),
	depthCallback func(*goex.Depth),
	tradeCallback func(*goex.Trade),
	klineCallback func(*goex.Kline, int)) {
	kcb := func(kline *goex.Kline, period goex.KlinePeriod) {
		klineCallback(kline, int(period))
	}
	o.ok.SetCallbacks(tickerCallback, depthCallback, tradeCallback, kcb, nil)
}

type OKExFutureV3Ws struct {
	ok *okex.OKExV3FuturesWs
}

func NewOKExFutureV3Ws() *OKExFutureV3Ws {
	return &OKExFutureV3Ws{
		ok: okex.NewOKExV3FuturesWs(nil),
	}
}

func (o *OKExFutureV3Ws) SubscribeTicker(pair goex.CurrencyPair, contract string) error {
	return o.ok.SubscribeTicker(pair, contract)
}

func (o *OKExFutureV3Ws) SubscribeDepth(pair goex.CurrencyPair, contract string, size int) error {
	return o.ok.SubscribeDepth(pair, contract, size)
}

func (o *OKExFutureV3Ws) SubscribeTrade(pair goex.CurrencyPair, contract string) error {
	return o.ok.SubscribeTrade(pair, contract)
}

func (o *OKExFutureV3Ws) SubscribeKline(pair goex.CurrencyPair, contract string, period int) error {
	return o.ok.SubscribeKline(pair, contract, period)
}

func (o *OKExFutureV3Ws) SetCallbacks(
	tickerCallback func(ticker *goex.FutureTicker),
	depthCallback func(*goex.Depth),
	tradeCallback func(*goex.Trade, string),
	klineCallback func(*goex.FutureKline, int)) {
	o.ok.SetCallbacks(tickerCallback, depthCallback, tradeCallback, klineCallback, nil)
}
