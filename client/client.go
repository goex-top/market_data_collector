package client

import (
	"github.com/goex-top/market_center"
	"github.com/goex-top/market_center_client"
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
	"os"
)

type Client struct {
	ExchangeName string
	CurrencyPair string
	contractType string
	spotApi      goex.API
	futureApi    goex.FutureRestAPI
	c            *market_center_client.Client
	isDirect     bool
	isSpot       bool
}

func NewClient(exchangeName, currencyPair, contractType string, c *market_center_client.Client) *Client {
	proxy := os.Getenv("HTTP_PROXY")
	var spotApi goex.API
	var futureApi goex.FutureRestAPI
	var direct = false
	var isSpot = !market_center.IsFutureExchange(exchangeName)
	if c == nil {
		if isSpot {
			spotApi = builder.NewAPIBuilder().HttpProxy(proxy).Build(exchangeName)
		} else {
			futureApi = builder.NewAPIBuilder().HttpProxy(proxy).BuildFuture(exchangeName)
		}
		direct = true
	}
	return &Client{
		ExchangeName: exchangeName,
		CurrencyPair: currencyPair,
		spotApi:      spotApi,
		futureApi:    futureApi,
		contractType: contractType,
		c:            c,
		isDirect:     direct,
		isSpot:       isSpot,
	}
}

func (c *Client) Close() {
	if c.c != nil {
		c.c.Close()
	}
}

func (c *Client) Name() string {
	return c.ExchangeName
}

func (c *Client) GetTicker() *goex.Ticker {
	var tick *goex.Ticker
	var err error
	if c.isDirect {
		if c.isSpot {
			tick, err = c.spotApi.GetTicker(goex.NewCurrencyPair2(c.CurrencyPair))
		} else {
			tick, err = c.futureApi.GetFutureTicker(goex.NewCurrencyPair2(c.CurrencyPair), c.contractType)
		}
	} else {
		if c.isSpot {
			tick, err = c.c.GetSpotTicker(c.ExchangeName, c.CurrencyPair)
		} else {
			tick, err = c.c.GetFutureTicker(c.ExchangeName, c.contractType, c.CurrencyPair)
		}
	}
	if err != nil {
		return nil
	}
	return tick
}

func (c *Client) GetDepth() *goex.Depth {
	var depth *goex.Depth
	var err error
	if c.isDirect {
		if c.isSpot {
			depth, err = c.spotApi.GetDepth(20, goex.NewCurrencyPair2(c.CurrencyPair))
		} else {
			depth, err = c.futureApi.GetFutureDepth(goex.NewCurrencyPair2(c.CurrencyPair), c.contractType, 20)
		}
	} else {
		if c.isSpot {
			depth, err = c.c.GetSpotDepth(c.ExchangeName, c.CurrencyPair)
		} else {
			depth, err = c.c.GetFutureDepth(c.ExchangeName, c.contractType, c.CurrencyPair)
		}
	}
	if err != nil {
		return nil
	}
	return depth
}

func (c *Client) GetKline() *goex.Kline {
	panic("not support yet")
}
